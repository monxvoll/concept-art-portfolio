import { config } from '../config.js';
import { initializeApp } from "https://www.gstatic.com/firebasejs/10.9.0/firebase-app.js";
import { getAuth, signInWithEmailAndPassword, signOut, onAuthStateChanged } from "https://www.gstatic.com/firebasejs/10.9.0/firebase-auth.js";

// Uses a placeholder config. The user should replace it with their own credentials.
const firebaseConfig = {
    apiKey: config.firebase?.apiKey || "YOUR_API_KEY",
    authDomain: config.firebase?.authDomain || "YOUR_AUTH_DOMAIN",
    projectId: config.firebase?.projectId || "YOUR_PROJECT_ID",
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);

export class Admin {
    constructor(modalElement) {
        this.modal = modalElement;
        
        if (!this.modal) return;
        
        this.closeBtn = this.modal.querySelector('#login-close');
        this.form = this.modal.querySelector('#login-form');
        this.overlay = this.modal.querySelector('#login-overlay');

        this.initAuthListener();

        if (this.form) {
            this.initEventListeners();
        }
    }

    initAuthListener() {
        onAuthStateChanged(auth, (user) => {
            if (user) {
                console.log("Admin authenticated with Firebase. Enabling Admin Mode...");
                document.body.classList.add('admin-mode');
                
                // Load third-party admin assets, then notify other components
                this.loadAdminAssets().then(() => {
                    document.dispatchEvent(new CustomEvent('adminModeActivated', { detail: { user } }));
                }).catch(err => {
                    console.error("Failed to load admin assets:", err);
                });
            } else {
                console.log("Admin logged out from Firebase.");
                document.body.classList.remove('admin-mode');
                document.dispatchEvent(new CustomEvent('adminModeDeactivated'));
            }
        });
    }

    initEventListeners() {
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.shiftKey && e.key.toLowerCase() === config.adminShortcut.toLowerCase()) {
                e.preventDefault();
                this.toggleModal();
            }
            if (e.key === 'Escape' && !this.modal.classList.contains('hidden')) {
                this.closeModal();
            }
        });

        if (this.closeBtn) {
            this.closeBtn.addEventListener('click', () => this.closeModal());
        }
        if (this.overlay) {
            this.overlay.addEventListener('click', () => this.closeModal());
        }

        this.form.addEventListener('submit', (e) => {
            e.preventDefault();
            this.authenticate();
        });
    }

    toggleModal() {
        if (auth.currentUser) {
            this.logout();
            return;
        }

        if (this.modal.classList.contains('hidden')) {
            this.modal.classList.remove('hidden');
            setTimeout(() => {
                const emailInput = this.modal.querySelector('#admin-email');
                if (emailInput) emailInput.focus();
            }, 100);
        } else {
            this.closeModal();
        }
    }

    async logout() {
        try {
            await signOut(auth);
            window.location.reload();
        } catch (error) {
            console.error("Logout error:", error);
        }
    }

    closeModal() {
        this.modal.classList.add('hidden');
    }

    async authenticate() {
        const emailInput = this.modal.querySelector('#admin-email');
        const passInput = this.modal.querySelector('#admin-pass');
        
        const email = emailInput ? emailInput.value : '';
        const password = passInput ? passInput.value : '';

        if (email.length > 0 && password.length > 0) {
            try {
                await signInWithEmailAndPassword(auth, email, password);
                this.closeModal();
                if (emailInput) emailInput.value = '';
                if (passInput) passInput.value = '';
            } catch (error) {
                console.error("Authentication failed:", error);
                alert("Login failed: " + error.message);
            }
        }
    }

    loadAdminAssets() {
        const assets = [
            { type: 'css', url: 'https://cdn.jsdelivr.net/npm/slim-select@2/dist/slimselect.css' },
            { type: 'js', url: 'https://cdn.jsdelivr.net/npm/slim-select@2/dist/slimselect.min.js' },
            { type: 'css', url: 'https://cdn.jsdelivr.net/npm/flatpickr/dist/flatpickr.min.css' },
            { type: 'js', url: 'https://cdn.jsdelivr.net/npm/flatpickr' }
        ];

        const loadPromises = assets.map(asset => {
            return new Promise((resolve, reject) => {
                if (asset.type === 'css') {
                    if (document.querySelector(`link[href="${asset.url}"]`)) return resolve();
                    const link = document.createElement('link');
                    link.rel = 'stylesheet';
                    link.href = asset.url;
                    link.onload = resolve;
                    link.onerror = reject;
                    const customStyles = document.querySelector('link[href="./assets/css/gallery.css"]');
                    if (customStyles) {
                        document.head.insertBefore(link, customStyles);
                    } else {
                        document.head.appendChild(link);
                    }
                } else if (asset.type === 'js') {
                    if (document.querySelector(`script[src="${asset.url}"]`)) return resolve();
                    const script = document.createElement('script');
                    script.src = asset.url;
                    script.onload = resolve;
                    script.onerror = reject;
                    document.head.appendChild(script);
                }
            });
        });

        return Promise.all(loadPromises);
    }
}
