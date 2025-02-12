package services

import (
	"encoding/json"
	"fmt"
	"github.com/sorawaslocked/CodeRivals/internal/entities"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
	"runtime"
	"time"
)

type ExecutionResult struct {
	Success     bool
	TestResults []*TestResult
	TimeMS      int64
	MemoryKB    int
	Error       string
}

type TestResult struct {
	Passed   bool
	Input    string
	Expected string
	Actual   string
	TimeMS   int64
}

type CodeExecutionService struct {
}

func NewCodeExecutionService() *CodeExecutionService {
	return &CodeExecutionService{}
}

func convertArgsAndExecute(i *interp.Interpreter, problem *entities.Problem, testCase *entities.ProblemTestCase) (interface{}, error) {
	// Parse the JSON inputs string into []interface{}
	var rawInputs []interface{}
	if err := json.Unmarshal([]byte(testCase.Input), &rawInputs); err != nil {
		return nil, fmt.Errorf("failed to parse inputs JSON: %v", err)
	}

	// Rest of the function remains the same
	convertedArgs := make([]reflect.Value, len(rawInputs))
	for idx, rawInput := range rawInputs {
		targetType := problem.InputTypes[idx]
		converted, err := convertToType(rawInput, targetType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %d: %v", idx, err)
		}
		convertedArgs[idx] = reflect.ValueOf(converted)
	}

	v, err := i.Eval(fmt.Sprintf("solution.%s", problem.MethodName))
	if err != nil {
		return nil, fmt.Errorf("failed to find solution function: %v", err)
	}

	fn := reflect.ValueOf(v.Interface())
	result := fn.Call(convertedArgs)

	if len(result) != 1 {
		return nil, fmt.Errorf("expected 1 return value, got %d", len(result))
	}

	return result[0].Interface(), nil
}

func convertToType(value interface{}, targetType string) (interface{}, error) {
	switch targetType {
	case "[]int":
		// Handle array/slice conversion
		if arr, ok := value.([]interface{}); ok {
			result := make([]int, len(arr))
			for i, v := range arr {
				if num, ok := v.(float64); ok { // JSON numbers are float64
					result[i] = int(num)
				} else {
					return nil, fmt.Errorf("invalid type for array element %v", v)
				}
			}
			return result, nil
		}
		return nil, fmt.Errorf("expected array, got %T", value)

	case "int":
		if num, ok := value.(float64); ok {
			return int(num), nil
		}
		return nil, fmt.Errorf("expected number, got %T", value)

	case "string":
		if str, ok := value.(string); ok {
			return str, nil
		}
		return nil, fmt.Errorf("expected string, got %T", value)

	case "[][]int":
		if arr, ok := value.([]interface{}); ok {
			result := make([][]int, len(arr))
			for i, v := range arr {
				if subArr, ok := v.([]interface{}); ok {
					result[i] = make([]int, len(subArr))
					for j, sv := range subArr {
						if num, ok := sv.(float64); ok {
							result[i][j] = int(num)
						} else {
							return nil, fmt.Errorf("invalid type for 2D array element %v", sv)
						}
					}
				} else {
					return nil, fmt.Errorf("invalid type for array element %v", v)
				}
			}
			return result, nil
		}
		return nil, fmt.Errorf("expected 2D array, got %T", value)

	// Add more type conversions as needed
	default:
		return nil, fmt.Errorf("unsupported type: %s", targetType)
	}
}

func (s *CodeExecutionService) ExecuteSolution(problem *entities.Problem, testCases []*entities.ProblemTestCase, submission *entities.ProblemSubmission) ExecutionResult {
	i := interp.New(interp.Options{})

	fullCode := fmt.Sprintf(`
        package solution
        
        %s
    `, submission.Code)

	if err := i.Use(stdlib.Symbols); err != nil {
		return ExecutionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	// Evaluate the code
	_, err := i.Eval(fullCode)
	if err != nil {
		return ExecutionResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	// Get initial memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startMemory := m.Alloc

	results := make([]*TestResult, len(testCases))
	success := true
	var testError string

	for idx, tc := range testCases {
		start := time.Now()

		result, err := convertArgsAndExecute(i, problem, tc)
		if err != nil {
			return ExecutionResult{Success: false, Error: err.Error()}
		}

		// Parse expected output from JSON string
		var expectedOutput interface{}
		if err := json.Unmarshal([]byte(tc.Output), &expectedOutput); err != nil {
			return ExecutionResult{Success: false, Error: fmt.Sprintf("failed to parse expected output: %v", err)}
		}

		expected, err := convertToType(expectedOutput, problem.OutputType)
		if err != nil {
			return ExecutionResult{Success: false, Error: fmt.Sprintf("failed to convert expected output: %v", err)}
		}

		passed := reflect.DeepEqual(result, expected)
		results[idx] = &TestResult{
			Passed:   passed,
			Input:    tc.Input,  // Already a JSON string
			Expected: tc.Output, // Already a JSON string
			Actual:   fmt.Sprintf("%v", result),
			TimeMS:   time.Since(start).Milliseconds(),
		}

		if !passed {
			success = false
			testError, _ = result.(string)

			break
		}
	}

	// Get final memory usage after all test cases
	runtime.ReadMemStats(&m)
	memoryUsed := (m.Alloc - startMemory) / 1024 // Convert to KB

	var totalTime int64

	for _, tr := range results {
		if tr != nil {
			totalTime += tr.TimeMS
		}
	}

	avgTime := totalTime / int64(len(results))

	return ExecutionResult{
		Success:     success,
		TestResults: results,
		MemoryKB:    int(memoryUsed),
		TimeMS:      avgTime,
		Error:       testError,
	}
}
