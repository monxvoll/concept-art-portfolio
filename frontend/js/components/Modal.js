// Artwork detail modal with project-based navigation, sharing, and admin editing
import { config } from '../config.js';

export class Modal {
    /**
     * @param {HTMLElement} modalElement - The modal's root DOM element
     * @param {Array} artworks - Full artwork list (used to group by project)
     */
    constructor(modalElement, artworks = []) {
        this.modal = modalElement;
        this.artworks = artworks;
        this.currentProjectGroup = [];
        this.currentIndex = 0;

        if (!this.modal) return;

        this.overlay = this.modal.querySelector('#modal-overlay');
        this.closeBtn = this.modal.querySelector('#modal-close');

        this.shareBtn = this.modal.querySelector('#modal-share');
        this.shareText = this.modal.querySelector('#modal-share-text');

        // Navigation arrows (visible when artwork belongs to a multi-piece project)
        this.prevBtn = this.modal.querySelector('#modal-prev');
        this.nextBtn = this.modal.querySelector('#modal-next');

        this.image = this.modal.querySelector('#modal-image');
        this.title = this.modal.querySelector('#modal-title');
        this.category = this.modal.querySelector('#modal-category');
        this.categorySelect = this.modal.querySelector('#modal-category-select');
        this.description = this.modal.querySelector('#modal-description');
        this.date = this.modal.querySelector('#modal-date');
        this.datePicker = this.modal.querySelector('#modal-date-picker');
        this.saveBtn = this.modal.querySelector('#modal-save-btn');

        if (this.datePicker) {
            const today = new Date().toISOString().split('T')[0];
            this.datePicker.setAttribute('max', today);
        }

        this.textPane = this.modal.querySelector('.modal__text-pane');

        // Fullscreen image view (opens when clicking the main image)
        this.fullscreenView = document.getElementById('fullscreen-view');
        this.fullscreenImage = document.getElementById('fullscreen-image');

        this.initEvents();
    }

    initEvents() {
        this.closeBtn.addEventListener('click', () => this.close());
        this.overlay.addEventListener('click', () => this.close());

        // Share: uses native Web Share API on mobile, clipboard fallback on desktop
        this.shareBtn.addEventListener('click', () => {
            const currentArt = this.currentProjectGroup[this.currentIndex];

            if (navigator.share) {
                this.shareBtn.classList.add('success');
                this.shareText.innerHTML = '&nbsp;:&nbsp&nbsp&nbsp&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;)';
                this.shareText.classList.remove('hidden');

                navigator.share({
                    title: currentArt.title,
                    text: `Check out ${currentArt.title} on ${config.name}'s portfolio`,
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

        // Keyboard navigation: Escape closes, arrows navigate project group
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
                // Don't capture arrows when the user is editing text fields
                if (document.activeElement.isContentEditable || document.activeElement.tagName === 'INPUT') return;

                if (e.key === 'ArrowLeft') this.prev();
                if (e.key === 'ArrowRight') this.next();
            }
        });

        // Enforce character limits on contenteditable title and description
        const enforceLimit = (element, maxLimit) => {
            element.addEventListener('keydown', (e) => {
                const allowedKeys = ['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown', 'Tab'];
                if (allowedKeys.includes(e.key) || e.ctrlKey || e.metaKey) return;

                if (element.innerText.length >= maxLimit) {
                    e.preventDefault();
                }
            });
            // Truncate pasted text to remaining character budget
            element.addEventListener('paste', (e) => {
                e.preventDefault();
                const text = (e.clipboardData || window.clipboardData).getData('text');
                const remaining = maxLimit - element.innerText.length;
                if (remaining > 0) {
                    const selection = window.getSelection();
                    if (!selection.rangeCount) return;
                    const range = selection.getRangeAt(0);
                    range.deleteContents();
                    const textNode = document.createTextNode(text.substring(0, remaining));
                    range.insertNode(textNode);
                    range.setStartAfter(textNode);
                    range.setEndAfter(textNode);
                    selection.removeAllRanges();
                    selection.addRange(range);
                }
            });
            element.addEventListener('input', () => {
                element.style.outline = '';
            });
        };

        enforceLimit(this.title, 100);
        enforceLimit(this.description, 2000);

        // Admin mode: make title/description editable and initialize third-party widgets
        document.addEventListener('adminModeActivated', () => {
            this.title.contentEditable = true;
            this.description.contentEditable = true;

            if (typeof SlimSelect !== 'undefined' && !this.slimSelect && this.categorySelect) {
                this.slimSelect = new SlimSelect({
                    select: this.categorySelect,
                    settings: { showSearch: false, closeOnSelect: true },
                    events: {
                        afterChange: () => {
                            if (document.activeElement) document.activeElement.blur();
                        }
                    }
                });

                // Sync the select to the currently displayed artwork
                if (this.currentProjectGroup.length > 0 && this.currentProjectGroup[this.currentIndex]) {
                    this.slimSelect.setSelected(this.currentProjectGroup[this.currentIndex].category);
                }
            }

            if (typeof flatpickr !== 'undefined' && !this.flatpickr) {
                this.flatpickr = flatpickr(this.datePicker, {
                    dateFormat: 'Y-m-d',
                    maxDate: 'today',
                    disableMobile: true,
                });
                if (this.datePicker.value) {
                    this.flatpickr.setDate(this.datePicker.value);
                }
            }
        });

        document.addEventListener('adminModeDeactivated', () => {
            this.title.contentEditable = false;
            this.description.contentEditable = false;
        });

        // Save edits: validate, update the artwork object, and close
        const saveEdits = () => {
            if (!this.currentProjectGroup || this.currentProjectGroup.length === 0) return;

            let hasError = false;
            const titleText = this.title.innerText.trim();
            const descText = this.description.innerText.trim();

            this.title.style.outline = '';
            this.description.style.outline = '';

            // Highlight empty required fields
            if (!titleText) {
                this.title.style.outline = '2px solid rgba(255, 71, 87, 0.9)';
                hasError = true;
            }
            if (!descText) {
                this.description.style.outline = '2px solid rgba(255, 71, 87, 0.9)';
                hasError = true;
            }

            if (hasError) return;

            const currentArt = this.currentProjectGroup[this.currentIndex];
            currentArt.title = titleText;
            currentArt.category = this.categorySelect ? this.categorySelect.value : this.category.innerText;
            this.category.innerText = currentArt.category;
            currentArt.description = descText;
            currentArt.artwork_date = this.datePicker.value;
            this.date.textContent = 'Created on ' + currentArt.artwork_date;

            const originalText = this.saveBtn.textContent;
            this.saveBtn.textContent = 'Saved!';
            this.saveBtn.classList.add('success');

            setTimeout(() => {
                this.saveBtn.textContent = originalText;
                this.saveBtn.classList.remove('success');
                this.close();
            }, 600);
        };

        this.saveBtn.addEventListener('click', saveEdits);

        this.categorySelect.addEventListener('change', () => {
            this.categorySelect.blur();
        });

        // Close the flatpickr calendar when the text pane scrolls to prevent it from floating
        if (this.textPane) {
            this.textPane.addEventListener('scroll', () => {
                if (this.flatpickr && this.flatpickr.isOpen) {
                    this.flatpickr.close();
                }
            });
        }
    }

    /**
     * Opens the modal for a given artwork.
     * Groups all artworks from the same project for arrow navigation.
     * @param {Object} artwork
     * @param {boolean} pushHistory - Whether to push a new browser history entry
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

        this.updateUI(pushHistory ? 'push' : 'none');

        if (isCurrentlyHidden) {
            document.body.classList.add('modal-open');
        }

        this.modal.classList.remove('hidden');
    }

    /**
     * Populate the modal's DOM elements with the current artwork data.
     * @param {'push'|'replace'|'none'} historyMode - 'push' adds a history entry (first open),
     *   'replace' swaps the current entry (navigating inside the open modal), 'none' leaves history untouched.
     */
    updateUI(historyMode = 'push') {
        const currentArt = this.currentProjectGroup[this.currentIndex];

        this.image.src = currentArt.image_url;
        this.image.alt = currentArt.title;
        this.title.innerText = currentArt.title;
        this.category.innerText = currentArt.category;
        if (this.categorySelect) {

            const existingCustom = this.categorySelect.querySelector('option[data-custom]');
            if (existingCustom) existingCustom.remove();

            // If the artwork's category doesn't match any option, add it temporarily
            const exists = Array.from(this.categorySelect.options).some(opt => opt.value === currentArt.category);
            if (!exists && currentArt.category) {
                const customOpt = document.createElement('option');
                customOpt.value = currentArt.category;
                customOpt.textContent = currentArt.category;
                customOpt.setAttribute('data-custom', 'true');
                this.categorySelect.insertBefore(customOpt, this.categorySelect.firstChild);
            }
            this.categorySelect.value = currentArt.category;
            if (this.slimSelect) {
                this.slimSelect.setSelected(currentArt.category);
            }
        }
        this.description.innerText = currentArt.description;
        this.date.textContent = 'Created on ' + currentArt.artwork_date;
        this.datePicker.placeholder = currentArt.artwork_date || 'Date';
        this.datePicker.value = currentArt.artwork_date;
        if (this.flatpickr) {
            this.flatpickr.setDate(currentArt.artwork_date);
        }

        const isAdmin = document.body.classList.contains('admin-mode');
        this.title.contentEditable = isAdmin;
        this.description.contentEditable = isAdmin;

        // Show navigation arrows only for multi-piece projects
        if (this.currentProjectGroup.length > 1) {
            this.prevBtn.classList.remove('hidden');
            this.nextBtn.classList.remove('hidden');
        } else {
            this.prevBtn.classList.add('hidden');
            this.nextBtn.classList.add('hidden');
        }

        if (historyMode !== 'none') {
            const newUrl = window.location.pathname + '?artwork=' + currentArt.slug;
            const state = { slug: currentArt.slug };
            // Replace while navigating inside the modal so "Back" returns straight to the gallery
            if (historyMode === 'replace') {
                window.history.replaceState(state, '', newUrl);
            } else {
                window.history.pushState(state, '', newUrl);
            }
        }

        if (this.textPane) {
            this.textPane.scrollTop = 0;
        }
    }

    next() {
        if (this.currentProjectGroup.length <= 1) return;
        this.currentIndex = (this.currentIndex + 1) % this.currentProjectGroup.length;
        this.updateUI('replace');
    }

    prev() {
        if (this.currentProjectGroup.length <= 1) return;
        this.currentIndex = (this.currentIndex - 1 + this.currentProjectGroup.length) % this.currentProjectGroup.length;
        this.updateUI('replace');
    }

    close(pushHistory = true) {
        if (this.slimSelect) this.slimSelect.close();
        if (this.flatpickr) this.flatpickr.close();

        // Skip transition when closing programmatically (e.g. from URL routing)
        if (!pushHistory) {
            this.modal.style.transition = 'none';
            const content = this.modal.querySelector('.modal__content');
            if (content) content.style.transition = 'none';
        }

        if (pushHistory) {
            window.history.pushState(null, '', window.location.pathname);
        }

        const cleanup = () => {
            document.body.classList.remove('modal-open');

            this.modal.classList.add('hidden');
            this.image.src = '';
            this.currentProjectGroup = [];
            this.currentIndex = 0;

            // Restore CSS transitions after a programmatic close
            if (!pushHistory) {
                requestAnimationFrame(() => {
                    this.modal.style.transition = '';
                    const content = this.modal.querySelector('.modal__content');
                    if (content) content.style.transition = '';
                });
            }
        };

        if (!pushHistory) {
            cleanup();
        } else {
            this.modal.classList.add('hidden');
            setTimeout(cleanup, 250);
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
