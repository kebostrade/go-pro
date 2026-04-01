// Gin Web App JavaScript
// This file is loaded on every page

(function() {
    'use strict';

    // Log page load
    console.log('Gin Web App loaded');

    // Check API status on page load
    document.addEventListener('DOMContentLoaded', function() {
        checkApiStatus();
    });

    // Fetch API health status
    function checkApiStatus() {
        var statusElement = document.getElementById('api-status');
        if (!statusElement) return;

        fetch('/api/v1/health')
            .then(function(response) {
                if (response.ok) {
                    return response.json();
                }
                throw new Error('API unavailable');
            })
            .then(function(data) {
                statusElement.textContent = 'API Status: ' + data.status + ' (v' + data.version + ')';
                statusElement.parentElement.classList.add('ok');
                statusElement.parentElement.classList.remove('error');
            })
            .catch(function(error) {
                statusElement.textContent = 'API Status: Error - ' + error.message;
                statusElement.parentElement.classList.add('error');
                statusElement.parentElement.classList.remove('ok');
            });
    }
})();
