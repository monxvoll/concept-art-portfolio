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
        this.textPane = document.querySelector('.modal__text-pane');

        // Fullscreen elements
        this.fullscreenView = document.getElementById('fullscreen-view');
        this.fullscreenImage = document.getElementById('fullscreen-image');


        this.initEvents();
    }

    initEvents() {
        this.closeBtn.addEventListener('click', () => this.close());
        this.overlay.addEventListener('click', () => this.close());

        this.shareBtn.addEventListener('click', () => {
            const currentArt = this.currentProjectGroup[this.currentIndex];

            if (navigator.share) {
                this.shareBtn.classList.add('success');
                this.shareText.innerHTML = '&nbsp;:&nbsp&nbsp&nbsp&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;)';
                this.shareText.classList.remove('hidden');

                navigator.share({
                    title: currentArt.title,
                    text: `Check out ${currentArt.title} on David Ponguta's portfolio`,
                    url: window.location.href
                }).then(() => {
                    this.shareBtn.classList.remove('success');
                    this.shareText.classList.add('hidden');
                }).catch(err => {
                    this.shareBtn.classList.remove('success');
                    this.shareText.classList.add('hidden');
                    console.error('Share cancelled or failed: ', err);
                });
            } else {
                navigator.clipboard.writeText(window.location.href).then(() => {
                    this.shareBtn.classList.add('success');
                    this.shareText.textContent = 'Copied!';
                    this.shareText.classList.remove('hidden');
                    setTimeout(() => {
                        this.shareBtn.classList.remove('success');
                        this.shareText.classList.add('hidden');
                    }, 1000);
                }).catch(err => {
                    console.error('Failed to copy link: ', err);
                });
            }
        });

        this.prevBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            this.prev();
        });

        this.nextBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            this.next();
        });

        this.image.addEventListener('click', () => this.openFullscreen());

        this.fullscreenView.addEventListener('click', (e) => {
            if (e.target === this.fullscreenView) this.closeFullscreen();
        });

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

        const isCurrentlyHidden = this.modal.classList.contains('hidden');

        this.updateUI(pushHistory);

        if (isCurrentlyHidden) {
            const computedStyle = window.getComputedStyle(document.body);
            const currentPaddingRight = parseFloat(computedStyle.paddingRight) || 0;
            this.scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;
            
            this.originalPaddingRight = document.body.style.paddingRight;
            this.originalBodyOverflow = document.body.style.overflow;

            document.body.style.paddingRight = `${currentPaddingRight + this.scrollbarWidth}px`;
            document.body.style.overflow = 'hidden';
            document.body.classList.add('modal-open');
        }

        this.modal.classList.remove('hidden');
    }

    /**
     * @param {boolean} pushHistory 
     */
    updateUI(pushHistory = true) {
        const currentArt = this.currentProjectGroup[this.currentIndex];

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

        if (pushHistory) {
            const newUrl = window.location.pathname + '?obra=' + currentArt.slug;
            window.history.pushState({ slug: currentArt.slug }, '', newUrl);
        }

        if (this.textPane) {
            this.textPane.scrollTop = 0;
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
        if (!pushHistory) {
            this.modal.style.transition = 'none';
            const content = this.modal.querySelector('.modal__content');
            if (content) content.style.transition = 'none';
        }

        this.modal.classList.add('hidden');

        if (pushHistory) {
            window.history.pushState(null, '', window.location.pathname);
        }

        const cleanup = () => {
            document.body.style.paddingRight = this.originalPaddingRight || '';
            document.body.style.overflow = this.originalBodyOverflow || '';
            document.body.classList.remove('modal-open');

            this.image.src = '';
            this.currentProjectGroup = [];
            this.currentIndex = 0;

            if (!pushHistory) {
                setTimeout(() => {
                    this.modal.style.transition = '';
                    const content = this.modal.querySelector('.modal__content');
                    if (content) content.style.transition = '';
                }, 50);
            }
        };

        if (!pushHistory) {
            cleanup();
        } else {
            setTimeout(cleanup, 400);
        }
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
