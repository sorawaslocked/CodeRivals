// web/static/js/topics.js
document.addEventListener('DOMContentLoaded', function() {
    const topicSearch = document.getElementById('topicSearch');
    const topicsList = document.getElementById('topicsList');
    const topicItems = topicsList.getElementsByClassName('topic-item');
    const topicCheckboxes = document.querySelectorAll('.topic-checkbox');

    // Topic search functionality
    topicSearch.addEventListener('input', function(e) {
        const searchTerm = e.target.value.toLowerCase();

        Array.from(topicItems).forEach(item => {
            const topicName = item.querySelector('.topic-name').textContent.toLowerCase();
            item.style.display = topicName.includes(searchTerm) ? '' : 'none';
        });
    });

    // Topic filter functionality
    function updateProblems() {
        const selectedTopics = Array.from(topicCheckboxes)
            .filter(cb => cb.checked)
            .map(cb => cb.value);

        // Get current URL and update search params
        const url = new URL(window.location);
        if (selectedTopics.length > 0) {
            url.searchParams.set('topics', selectedTopics.join(','));
        } else {
            url.searchParams.delete('topics');
        }

        // Preserve page parameter if it exists
        const page = url.searchParams.get('page');
        if (!page || page === '1') {
            url.searchParams.delete('page');
        }

        // Update URL without reloading
        window.history.pushState({}, '', url);

        // Fetch and update problems list
        fetch(url)
            .then(response => response.text())
            .then(html => {
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');
                const newProblemsList = doc.querySelector('.problems-list');
                const newPagination = doc.querySelector('.pagination');

                document.querySelector('.problems-list').innerHTML = newProblemsList.innerHTML;
                const paginationContainer = document.querySelector('.pagination');
                if (paginationContainer) {
                    paginationContainer.innerHTML = newPagination ? newPagination.innerHTML : '';
                }
            });
    }

    // Add event listeners to checkboxes
    topicCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', updateProblems);
    });

    // Initialize selected topics from URL
    const url = new URL(window.location);
    const topicsParam = url.searchParams.get('topics');
    if (topicsParam) {
        const selectedTopics = topicsParam.split(',');
        topicCheckboxes.forEach(checkbox => {
            checkbox.checked = selectedTopics.includes(checkbox.value);
        });
    }
});