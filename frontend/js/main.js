import { artworks } from './services/mockData.js';
import { Collage } from './components/Collage.js';
import { Modal } from './components/Modal.js';

document.addEventListener('DOMContentLoaded', () => {
    const collageContainer = document.getElementById('collage-container');
    const modal = new Modal(artworks);
    
    const collage = new Collage(collageContainer, artworks, (artwork) => {
        modal.open(artwork);
    });

    collage.render();

    // Check URL for deep linking
    const urlParams = new URLSearchParams(window.location.search);
    const obraSlug = urlParams.get('obra');
    if (obraSlug) {
        const targetArt = artworks.find(a => a.slug === obraSlug);
        if (targetArt) {
            // Open without pushing state again since it's already in the URL
            modal.open(targetArt, false);
        }
    }

    // Handle browser Back/Forward buttons (popstate)
    window.addEventListener('popstate', () => {
        const params = new URLSearchParams(window.location.search);
        const slug = params.get('obra');

        if (slug) {
            const targetArt = artworks.find(a => a.slug === slug);
            if (targetArt) {
                modal.open(targetArt, false);
            }
        } else {
            modal.close(false);
        }
    });

    // --- Resume Modal Logic ---
    const resumeBtn = document.getElementById('nav-resume');
    const resumeModal = document.getElementById('resume-modal');
    const resumeClose = document.getElementById('resume-close');
    const resumeOverlay = document.getElementById('resume-overlay');

    const resumeContent = document.getElementById('resume-content');

    const openResume = () => {
        if (resumeContent) {
            resumeContent.scrollTop = 0;
        }
        resumeModal.classList.remove('hidden');
        document.body.style.overflow = 'hidden'; // Prevent scrolling
    };

    const closeResume = () => {
        resumeModal.classList.add('hidden');
        document.body.style.overflow = '';
    };

    resumeBtn.addEventListener('click', openResume);
    resumeClose.addEventListener('click', closeResume);
    resumeOverlay.addEventListener('click', closeResume);

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !resumeModal.classList.contains('hidden')) {
            closeResume();
        }
    });

    // --- Back to top button logic ---
    const backToTopBtn = document.getElementById('back-to-top');

    window.addEventListener('scroll', () => {
        if (window.scrollY > 300) {
            backToTopBtn.classList.add('visible');
        } else {
            backToTopBtn.classList.remove('visible');
        }
    }, { passive: true });

    // Use IntersectionObserver to seamlessly pin the button above the footer
    const footer = document.querySelector('.footer');
    if (footer && backToTopBtn) {
        document.body.style.position = 'relative'; // Ensure body is the relative container
        
        const footerObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const footerHeight = entry.boundingClientRect.height;
                    // Switch to absolute positioning tied to the bottom of the body
                    backToTopBtn.style.position = 'absolute';
                    backToTopBtn.style.bottom = `calc(${footerHeight}px + 2rem)`;
                    backToTopBtn.style.transform = 'none'; // reset any residual transform
                } else {
                    // Reset to normal fixed positioning
                    backToTopBtn.style.position = '';
                    backToTopBtn.style.bottom = '';
                    backToTopBtn.style.transform = '';
                }
            });
        }, { 
            root: null,
            threshold: 0,
            // Small positive margin so it triggers just before hitting the footer
            rootMargin: '50px 0px 0px 0px' 
        });

        footerObserver.observe(footer);
    }

    backToTopBtn.addEventListener('click', (e) => {
        e.preventDefault();
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    });

    // Force trigger during momentum scroll on mobile
    backToTopBtn.addEventListener('touchstart', (e) => {
        e.preventDefault(); // Prevents the OS from swallowing the tap to stop scroll
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    }, { passive: false });
});
