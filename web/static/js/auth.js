document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                })
            });

            const data = await response.json();

            if (response.ok) {
                // Successful login
                window.location.href = '/dashboard';
            } else {
                // Handle error
                const errorDiv = document.createElement('div');
                errorDiv.className = 'error-message';
                errorDiv.textContent = data.error || 'Login failed. Please try again.';

                const existingError = document.querySelector('.error-message');
                if (existingError) {
                    existingError.remove();
                }

                loginForm.insertBefore(errorDiv, loginForm.firstChild);
            }
        } catch (error) {
            console.error('Login error:', error);
        }
    });
});