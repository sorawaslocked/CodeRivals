package entities

type LearningResource struct {
	Title        string
	Description  string
	ResourceLink string
	Difficulty   string
	Complexity   string
	Author       string
}

type ProblemSolvingStrategy struct {
	Title       string
	Description string
}

type LearningMaterials struct {
	DataStructures           []LearningResource
	Algorithms               []LearningResource
	Books                    []LearningResource
	OnlineResources          []LearningResource
	ProblemSolvingStrategies []ProblemSolvingStrategy
}

func GetLearningMaterials() LearningMaterials {
	return LearningMaterials{
		DataStructures: []LearningResource{
			{
				Title:        "Arrays and Strings",
				Description:  "Master fundamental data structures for storing and manipulating collections of elements.",
				ResourceLink: "https://www.geeksforgeeks.org/array-data-structure/",
				Difficulty:   "Beginner",
			},
			{
				Title:        "Linked Lists",
				Description:  "Deep dive into dynamic data structures with nodes and pointers.",
				ResourceLink: "https://www.topcoder.com/community/competitive-programming/tutorials/data-structures/",
				Difficulty:   "Intermediate",
			},
		},
		Algorithms: []LearningResource{
			{
				Title:        "Sorting Algorithms",
				Description:  "Comprehensive guide to essential sorting techniques like Bubble Sort, Quick Sort, and Merge Sort.",
				ResourceLink: "https://www.programiz.com/dsa/sorting-algorithm",
				Complexity:   "O(n log n)",
			},
			{
				Title:        "Graph Traversal",
				Description:  "Learn depth-first and breadth-first search algorithms with practical implementations.",
				ResourceLink: "https://www.geeksforgeeks.org/graph-data-structure-and-algorithms/",
				Complexity:   "O(V + E)",
			},
		},
		Books: []LearningResource{
			{
				Title:        "Cracking the Coding Interview",
				Description:  "The ultimate guide for coding interview preparation with 189 programming questions and solutions.",
				ResourceLink: "https://www.amazon.com/Cracking-Coding-Interview-Programming-Questions/dp/0984782850",
				Author:       "Gayle Laakmann McDowell",
			},
			{
				Title:        "Introduction to Algorithms",
				Description:  "A comprehensive textbook covering a wide range of algorithms with in-depth explanations.",
				ResourceLink: "https://www.amazon.com/Introduction-Algorithms-3rd-MIT-Press/dp/0262033844",
				Author:       "Thomas H. Cormen, et al.",
			},
			{
				Title:        "The Algorithm Design Manual",
				Description:  "Practical guide to designing and implementing efficient algorithms with real-world examples.",
				ResourceLink: "https://www.amazon.com/Algorithm-Design-Manual-Steven-Skiena/dp/1849967202",
				Author:       "Steven S. Skiena",
			},
			{
				Title:        "Grokking Algorithms: An Illustrated Guide for Programmers and Other Curious People",
				Description:  "A beginner-friendly, illustrated guide to common algorithms like sorting, searching, and data compression, with step-by-step explanations and Python code samples.",
				ResourceLink: "https://www.amazon.com/Grokking-Algorithms-illustrated-programmers-curious/dp/1617292230",
				Author:       "Aditya Bhargava",
			},
		},
		OnlineResources: []LearningResource{
			{
				Title:        "GeeksforGeeks",
				Description:  "Platform with 2000+ coding problems and detailed solutions across various difficulty levels.",
				ResourceLink: "https://www.geeksforgeeks.org/",
			},
			{
				Title:        "Khan Academy - Algorithms",
				Description:  "beginner-friendly course that introduces fundamental algorithms using interactive lessons and visual explanations. It covers sorting, searching, recursion, and graph algorithms with step-by-step walkthroughs, making it great for those new to the topic. The platform’s engaging approach helps learners grasp complex concepts intuitively.",
				ResourceLink: "https://www.khanacademy.org/computing/computer-science/algorithms",
			},
		},
	}
}
