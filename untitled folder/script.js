// Format number with Swedish locale (sv-SE) - space as thousand separator
function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
}

// Clamp value between min and max
function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max);
}

const API_URL = '/api/state';

// Update overlay from server
async function updateOverlay() {
    try {
        const response = await fetch(API_URL, {
            cache: 'no-store'
        });
        
        if (!response.ok) {
            throw new Error('Failed to fetch state');
        }
        
        const state = await response.json();
        
        // Use values from server with defaults
        const label = state.label || 'LIKE GOAL';
        const likes = parseInt(state.likes || 0, 10);
        const goal = parseInt(state.goal || 1000000, 10);
        const rotationDuration = parseFloat(state.rotationDuration || 4);
        const rotationSpeed = parseFloat(state.rotationSpeed || 1);

        // Update label
        const labelElement = document.getElementById('likeGoalLabel');
        if (labelElement) {
            labelElement.textContent = label.toUpperCase();
        }

        // Update current likes
        const currentLikesElement = document.getElementById('currentLikes');
        if (currentLikesElement) {
            currentLikesElement.textContent = formatNumber(likes);
        }

        // Update goal
        const goalLikesElement = document.getElementById('goalLikes');
        if (goalLikesElement) {
            goalLikesElement.textContent = formatNumber(goal);
        }

        // Calculate and update progress bar
        const progressPercent = clamp((likes / goal) * 100, 0, 100);
        const progressBarElement = document.getElementById('progressBar');
        const progressBarContainer = progressBarElement?.parentElement;
        const pokeballIcon = document.querySelector('.pokeball-icon');
        
        if (progressBarElement && progressBarContainer) {
            progressBarElement.style.width = progressPercent + '%';
            
            // Position Pokeball icon at the tip of the progress bar
            if (pokeballIcon) {
                // Use requestAnimationFrame to ensure layout has updated
                requestAnimationFrame(() => {
                    const progressBarRect = progressBarElement.getBoundingClientRect();
                    const containerRect = progressBarContainer.getBoundingClientRect();
                    const pokeballWidth = 55;
                    // Calculate position relative to container (no padding now)
                    const progressBarRight = progressBarRect.right - containerRect.left;
                    const pokeballLeft = progressBarRight - (pokeballWidth / 2);
                    pokeballIcon.style.left = Math.max(0, pokeballLeft) + 'px';
                });
            }
        }

        // Update logo rotation - duration controls cycle interval, speed controls rotation timing
        const logoElement = document.querySelector('.pokexclusive-logo');
        if (logoElement) {
            // Calculate the keyframe percentage where rotation should start
            // If rotationSpeed is 1s and rotationDuration is 4s, rotation takes 25% of the cycle
            // So rotation should start at 75% (100% - 25%)
            const rotationPercent = Math.min(Math.max(rotationSpeed / rotationDuration, 0.05), 0.95); // Clamp between 5% and 95%
            const stillPercent = (1 - rotationPercent) * 100;
            
            // Use rotationDuration for the total cycle time (interval)
            logoElement.style.animationDuration = rotationDuration + 's';
            
            // Dynamically update keyframes by creating/updating a style element
            let styleElement = document.getElementById('dynamic-rotation-keyframes');
            if (!styleElement) {
                styleElement = document.createElement('style');
                styleElement.id = 'dynamic-rotation-keyframes';
                document.head.appendChild(styleElement);
            }
            
            styleElement.textContent = `
                @keyframes logoFullRotation {
                    0%, ${stillPercent}% {
                        transform: perspective(1000px) rotateY(0deg);
                        filter: 
                            drop-shadow(0 0 15px rgba(183, 148, 246, 0.8))
                            drop-shadow(0 0 30px rgba(183, 148, 246, 0.6));
                    }
                    100% {
                        transform: perspective(1000px) rotateY(360deg);
                        filter: 
                            drop-shadow(0 0 30px rgba(183, 148, 246, 1))
                            drop-shadow(0 0 60px rgba(183, 148, 246, 0.9))
                            drop-shadow(0 0 90px rgba(183, 148, 246, 0.7));
                    }
                }
            `;
        }
    } catch (error) {
        console.error('Error updating overlay:', error);
    }
}

// Initialize overlay
updateOverlay();

// Poll server every 300ms for updates
setInterval(updateOverlay, 300);
