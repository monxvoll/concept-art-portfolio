// Entry point - initializes all components and handles page-level routing
import { fetchArtworks, deleteArtwork, createArtwork, fetchResumeURL, uploadResume } from './services/api.js';
import { Collage } from './components/Collage.js';
import { Modal } from './components/Modal.js';
import { Admin, auth } from './components/AdminAuth.js';
import { AddArtwork } from './components/AddArtwork.js';
import { config } from './config.js';

let artworks = [];

document.addEventListener('DOMContentLoaded', async () => {

    // Inject config values (name, role, email, etc.) into the DOM
    const initConfig = () => {
        const setEl = (id, text) => {
            const el = document.getElementById(id);
            if (el) el.textContent = text;
        };
        const setAttr = (id, attr, val) => {
            const el = document.getElementById(id);
            if (el) el.setAttribute(attr, val);
        };

        setEl('config-name', config.name);
        setEl('config-role', config.role);
        setEl('config-resume-name', config.name);
        setEl('config-resume-role', config.role);
        setEl('config-login-subtitle', config.name);
        setEl('config-copyright', `© ${config.year} ${config.name}. All rights reserved.`);

        setAttr('config-contact', 'href', `mailto:${config.email}`);
        setAttr('config-footer-contact', 'href', `mailto:${config.email}`);
        setEl('config-footer-contact', `Drop me a line at ${config.email}`);
        document.title = `${config.role} Portfolio - ${config.name}`;

        // Populate category <select> dropdowns from config
        const categories = config.categories || [];
        const popSelect = (id) => {
            const sel = document.getElementById(id);
            if (!sel) return;
            const firstOpt = sel.options[0];
            sel.innerHTML = '';
            sel.appendChild(firstOpt);
            categories.forEach(cat => {
                const opt = document.createElement('option');
                opt.value = cat;
                opt.textContent = cat;
                sel.appendChild(opt);
            });
        };
        popSelect('add-artwork-category');
        popSelect('modal-category-select');
    };
    initConfig();

    try {
        artworks = await fetchArtworks();
    } catch (e) {
        console.error("Failed to fetch artworks:", e);
        // Fallback to empty array if backend is down
        artworks = [];
    }

    // Prevent browser from restoring scroll position on reload
    if ('scrollRestoration' in history) {
        history.scrollRestoration = 'manual';
    }

    const byDateDesc = (a, b) => new Date(b.artwork_date) - new Date(a.artwork_date);

    // Read the artwork slug from the URL query string (?artwork=slug)
    const getUrlSlug = () => new URLSearchParams(window.location.search).get('artwork');

    // Core DOM references
    const collageContainer = document.getElementById('collage-container');
    const artworkModalEl = document.getElementById('artwork-modal');
    const loginModalEl = document.getElementById('login-modal');
    const addArtworkModalEl = document.getElementById('add-artwork-modal');

    // Initialize components
    const modal = new Modal(artworkModalEl, artworks);
    const admin = new Admin(loginModalEl);
    
    const addArtworkModal = new AddArtwork(addArtworkModalEl, async (formData) => {
        try {
            if (!auth.currentUser) throw new Error("Not logged in");
            const token = await auth.currentUser.getIdToken();
            const newArtwork = await createArtwork(formData, token);
            artworks.push(newArtwork);
            artworks.sort(byDateDesc);
            collage.render();
        } catch (e) {
            console.error("Failed to add artwork", e);
            alert("Failed to add artwork");
        }
    });

    // Sort artworks newest-first before initial render
    artworks.sort(byDateDesc);

    const collage = new Collage(collageContainer, artworks, (artwork) => {
        modal.open(artwork);
    });

    collage.render();

    // -- Resume Modal --

    const resumeBtn = document.getElementById('nav-resume');
    const resumeModal = document.getElementById('resume-modal');
    const resumeClose = document.getElementById('resume-close');
    const resumeOverlay = document.getElementById('resume-overlay');
    const resumeContent = document.getElementById('resume-content');
    const resumeDownloadBtn = document.getElementById('resume-download-btn');
    const resumeUploadFile = document.getElementById('resume-upload-file');
    const resumeUploadBtn = document.getElementById('resume-upload-btn');

    // Load initial resume URL
    fetchResumeURL().then(data => {
        if (data.url) resumeDownloadBtn.href = data.url;
    }).catch(e => console.error("Failed to load resume URL", e));

    if (resumeUploadBtn && resumeUploadFile) {
        resumeUploadBtn.addEventListener('click', () => {
            resumeUploadFile.click();
        });

        resumeUploadFile.addEventListener('change', async (e) => {
            const file = e.target.files[0];
            if (!file) return;

            try {
                if (!auth.currentUser) throw new Error("Not logged in");
                const token = await auth.currentUser.getIdToken();
                
                resumeUploadBtn.textContent = 'Uploading...';
                resumeUploadBtn.disabled = true;

                const formData = new FormData();
                formData.append('resume', file);
                
                const data = await uploadResume(formData, token);
                if (data.url) {
                    resumeDownloadBtn.href = data.url;
                    alert("Resume updated!");
                }
            } catch (err) {
                console.error("Resume upload failed", err);
                alert("Upload failed.");
            } finally {
                resumeUploadBtn.textContent = 'Upload New PDF';
                resumeUploadBtn.disabled = false;
                resumeUploadFile.value = '';
            }
        });
    }

    const openResume = (pushHistory = true) => {
        resumeContent.scrollTop = 0;
        resumeModal.classList.remove('hidden');
        document.body.classList.add('modal-open');
        if (pushHistory) {
            window.history.pushState({ resume: true }, '', window.location.pathname + '#resume');
        }
    };

    const closeResume = (pushHistory = true) => {
        resumeModal.classList.add('hidden');
        document.body.classList.remove('modal-open');
        if (pushHistory) {
            window.history.pushState(null, '', window.location.pathname);
        }
    };

    resumeBtn.addEventListener('click', () => openResume(true));
    resumeClose.addEventListener('click', () => closeResume(true));
    resumeOverlay.addEventListener('click', () => closeResume(true));

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !resumeModal.classList.contains('hidden')) {
            closeResume(true);
        }
    });

    // -- URL Routing --
    // Syncs modal state (artwork / resume) with the browser URL on load and popstate

    const handleRoute = () => {
        const slug = getUrlSlug();
        const isResume = window.location.hash === '#resume';

        if (isResume) {
            openResume(false);
        } else if (!resumeModal.classList.contains('hidden')) {
            closeResume(false);
        }

        if (slug && !isResume) {
            const targetArt = artworks.find(a => a.slug === slug);
            if (targetArt) {
                modal.open(targetArt, false);
            } else {
                modal.close(false);
            }
        } else {
            modal.close(false);
        }
    };

    handleRoute();
    window.addEventListener('popstate', handleRoute);

    // -- Back-to-top button --

    const backToTopBtn = document.getElementById('back-to-top');

    window.addEventListener('scroll', () => {
        if (window.scrollY > 300) {
            backToTopBtn.classList.add('visible');
        } else {
            backToTopBtn.classList.remove('visible');
        }
    }, { passive: true });

    backToTopBtn.addEventListener('click', (e) => {
        e.preventDefault();
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    });

    backToTopBtn.addEventListener('touchstart', (e) => {
        e.preventDefault();
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    }, { passive: false });

    // -- Delete artwork confirmation (mock to API) --

    const confirmModal = document.getElementById('confirm-modal');
    const confirmOverlay = document.getElementById('confirm-overlay');
    const confirmCancel = document.getElementById('confirm-cancel');
    const confirmDelete = document.getElementById('confirm-delete');
    const confirmArtTitle = document.getElementById('confirm-art-title');
    let artworkToDelete = null;

    document.addEventListener('requestDeleteArtwork', (e) => {
        artworkToDelete = e.detail;
        confirmArtTitle.textContent = artworkToDelete.title;
        confirmModal.classList.remove('hidden');
    });

    const closeConfirmModal = () => {
        confirmModal.classList.add('hidden');
        artworkToDelete = null;
        if (document.activeElement) document.activeElement.blur();
    };

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !confirmModal.classList.contains('hidden')) {
            closeConfirmModal();
        }
    });

    if (confirmOverlay) confirmOverlay.addEventListener('click', closeConfirmModal);
    if (confirmCancel) confirmCancel.addEventListener('click', closeConfirmModal);

    if (confirmDelete) confirmDelete.addEventListener('click', async () => {
        if (artworkToDelete) {
            try {
                if (!auth.currentUser) throw new Error("Not logged in");
                const token = await auth.currentUser.getIdToken();
                await deleteArtwork(artworkToDelete.id, token);
                
                const index = artworks.findIndex(a => a.id === artworkToDelete.id);
                if (index > -1) {
                    artworks.splice(index, 1);
                    collage.render();
                }
            } catch (e) {
                console.error("Failed to delete artwork", e);
                alert("Failed to delete artwork");
            }
            closeConfirmModal();
        }
    });

    // -- Add Artwork button (admin only, visibility handled via CSS) --

    const addArtworkBtn = document.getElementById('add-artwork-btn');
    if (addArtworkBtn) {
        addArtworkBtn.addEventListener('click', () => {
            addArtworkModal.open();
        });
    }
});
