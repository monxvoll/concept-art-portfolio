export class Modal {
    /**
     * @param {Array} artworks - Full database of artworks to enable project grouping
     */
    constructor(artworks = []) {
        this.artworks = artworks;
        this.currentProjectGroup = [];
        this.currentIndex = 0;

        // Base modal elements
        this.modal = document.getElementById('artwork-modal');
        this.overlay = document.getElementById('modal-overlay');
        this.closeBtn = document.getElementById('modal-close');

        // Share elements
        this.shareBtn = document.getElementById('modal-share');
        this.shareText = document.getElementById('modal-share-text');

        // Navigation arrows
        this.prevBtn = document.getElementById('modal-prev');
        this.nextBtn = document.getElementById('modal-next');

        // Data elements
        this.image = document.getElementById('modal-image');
        this.title = document.getElementById('modal-title');
        this.category = document.getElementById('modal-category');
        this.description = document.getElementById('modal-description');

        // Fullscreen elements
        this.fullscreenView = document.getElementById('fullscreen-view');
        this.fullscreenImage = document.getElementById('fullscreen-image');
        this.fullscreenCloseBtn = document.getElementById('fullscreen-close');

        this.initEvents();
    }

    initEvents() {
        // --- Base Modal Events ---

        // Close modal when X is clicked
        this.closeBtn.addEventListener('click', () => this.close());

        // Close modal when clicking outside the content (on the overlay)
        this.overlay.addEventListener('click', () => this.close());

        // Share button click
        this.shareBtn.addEventListener('click', () => {
            navigator.clipboard.writeText(window.location.href).then(() => {
                this.shareBtn.classList.add('success');
                this.shareText.classList.remove('hidden');

                // Hide success state after 1 second
                setTimeout(() => {
                    this.shareBtn.classList.remove('success');
                    this.shareText.classList.add('hidden');
                }, 1000);
            }).catch(err => {
                console.error('Failed to copy link: ', err);
            });
        });

        // Arrow clicks
        this.prevBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            this.prev();
        });

        this.nextBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            this.next();
        });

        // --- Fullscreen Events ---
        this.image.addEventListener('click', () => this.openFullscreen());
        this.fullscreenCloseBtn.addEventListener('click', () => this.closeFullscreen());
        this.fullscreenView.addEventListener('click', (e) => {
            if (e.target === this.fullscreenView) this.closeFullscreen();
        });

        // --- Keyboard Events ---
        document.addEventListener('keydown', (e) => {
            if (this.modal.classList.contains('hidden')) return;

            if (e.key === 'Escape') {
                if (!this.fullscreenView.classList.contains('hidden')) {
                    this.closeFullscreen();
                    return;
                }
                if (!this.modal.classList.contains('hidden')) {
                    this.close();
                }
            }

            if (!this.modal.classList.contains('hidden') && this.fullscreenView.classList.contains('hidden')) {
                if (e.key === 'ArrowLeft') this.prev();
                if (e.key === 'ArrowRight') this.next();
            }
        });
    }

    /**
     * @param {Object} artwork 
     * @param {boolean} pushHistory - Whether to push a new state to browser history
     */
    open(artwork, pushHistory = true) {
        if (artwork.project) {
            this.currentProjectGroup = this.artworks.filter(a => a.project === artwork.project);
            this.currentIndex = this.currentProjectGroup.findIndex(a => a.id === artwork.id);
        } else {
            this.currentProjectGroup = [artwork];
            this.currentIndex = 0;
        }

        this.updateUI(pushHistory);
        document.body.style.overflow = 'hidden';
        this.modal.classList.remove('hidden');
    }

    /**
     * @param {boolean} pushHistory 
     */
    updateUI(pushHistory = true) {
        const currentArt = this.currentProjectGroup[this.currentIndex];

        // Draw the artwork on the modal 
        this.image.src = currentArt.image_url;
        this.image.alt = currentArt.title;
        this.title.textContent = currentArt.title;
        this.category.textContent = currentArt.category;
        this.description.textContent = currentArt.description;

        if (this.currentProjectGroup.length > 1) {
            this.prevBtn.classList.remove('hidden');
            this.nextBtn.classList.remove('hidden');
        } else {
            this.prevBtn.classList.add('hidden');
            this.nextBtn.classList.add('hidden');
        }

        // Update browser URL
        if (pushHistory) {
            const newUrl = window.location.pathname + '?obra=' + currentArt.slug;
            window.history.pushState({ slug: currentArt.slug }, '', newUrl);
        }
    }

    next() {
        if (this.currentProjectGroup.length <= 1) return;
        this.currentIndex = (this.currentIndex + 1) % this.currentProjectGroup.length;
        this.updateUI();
    }

    prev() {
        if (this.currentProjectGroup.length <= 1) return;
        this.currentIndex = (this.currentIndex - 1 + this.currentProjectGroup.length) % this.currentProjectGroup.length;
        this.updateUI();
    }

    close(pushHistory = true) {
        this.modal.classList.add('hidden');
        document.body.style.overflow = 'auto';

        // Clean up URL if we are closing normally (not via back button)
        if (pushHistory) {
            window.history.pushState(null, '', window.location.pathname);
        }

        setTimeout(() => {
            this.image.src = '';
            this.currentProjectGroup = [];
            this.currentIndex = 0;
        }, 400);
    }

    openFullscreen() {
        this.fullscreenImage.src = this.image.src;
        this.fullscreenImage.alt = this.image.alt;
        this.fullscreenView.classList.remove('hidden');
    }

    closeFullscreen() {
        this.fullscreenView.classList.add('hidden');
        setTimeout(() => {
            this.fullscreenImage.src = '';
        }, 300);
    }
}
