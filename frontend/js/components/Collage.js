export class Collage {
    /**
     * @param {HTMLElement} container - The DOM element where the collage will be injected
     * @param {Array} artworks - Array of artworks (ARTWORK model)
     * @param {Function} onArtworkClick - Callback executed when an artwork is clicked
     */
    constructor(container, artworks, onArtworkClick) {
        this.container = container;
        this.artworks = artworks;
        this.onArtworkClick = onArtworkClick;
    }

    render() {
        this.container.innerHTML = ''; // Clear the container

        // Setup observer for entrance animations
        let appearDelay = 0;
        let appearTimeout = null;
        const appearObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    setTimeout(() => {
                        entry.target.classList.add('is-visible');
                    }, appearDelay);
                    appearDelay += 100; // Stagger effect

                    clearTimeout(appearTimeout);
                    appearTimeout = setTimeout(() => {
                        appearDelay = 0;
                    }, 200);

                    appearObserver.unobserve(entry.target);
                }
            });
        }, { rootMargin: '50px 0px', threshold: 0.05 });

        // Setup observer for mobile "hover" effect (active state on scroll)
        const hasHover = window.matchMedia('(hover: hover)').matches;
        let observer = null;

        if (!hasHover) {
            const observerOptions = {
                root: null,
                rootMargin: '-49.5% 0px -49.5% 0px', // Trigger when in the middle 20% of viewport
                threshold: 0
            };

            observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        entry.target.classList.add('is-active');
                    } else {
                        entry.target.classList.remove('is-active');
                    }
                });
            }, observerOptions);

            let wasAtBottom = false;
            window.addEventListener('scroll', () => {
                const items = document.querySelectorAll('.collage__item');
                if (items.length === 0) return;
                
                const lastItem = items[items.length - 1];
                // Increased threshold to 150px for more reliable detection on mobile browsers
                const isAtBottom = (window.innerHeight + window.scrollY) >= document.body.offsetHeight - 150;
                
                if (isAtBottom) {
                    if (!wasAtBottom) {
                        items.forEach(item => {
                            if (item !== lastItem) item.classList.remove('is-active');
                        });
                        lastItem.classList.add('is-active');
                        wasAtBottom = true;
                    }
                } else {
                    if (wasAtBottom) {
                        lastItem.classList.remove('is-active');
                        // We just left the bottom, manually check and restore the active state
                        // for whichever item is currently in the center of the screen
                        const center = window.innerHeight / 2;
                        items.forEach(item => {
                            const r = item.getBoundingClientRect();
                            if (r.top <= center && r.bottom >= center) {
                                item.classList.add('is-active');
                            }
                        });
                        wasAtBottom = false;
                    }
                }
            }, { passive: true });
        }

        this.artworks.forEach((art) => {
            const item = document.createElement('div');
            item.className = 'collage__item';

            const img = document.createElement('img');
            img.src = art.image_url;
            img.alt = art.title;
            img.className = 'collage__image';
            img.loading = 'lazy'; // Native lazy loading

            // Add click listener
            item.addEventListener('click', () => {
                if (this.onArtworkClick) {
                    this.onArtworkClick(art);
                }
            });

            // Create overlay container
            const overlay = document.createElement('div');
            overlay.className = 'collage__overlay';

            // Create title
            const title = document.createElement('h3');
            title.className = 'collage__title';
            title.textContent = art.title;

            // Assemble the puzzle
            overlay.appendChild(title);
            item.appendChild(img);
            item.appendChild(overlay);
            this.container.appendChild(item);

            // Observe the item for entrance animation
            appearObserver.observe(item);

            if (observer) {
                observer.observe(item);
            }
        });
    }
}
