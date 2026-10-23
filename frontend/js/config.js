// Portfolio configuration - edit these values to personalize the site
export const config = {
    name: "David Ponguta",
    role: "Concept Artist",
    email: "ponguart@gmail.com",
    year: "2026",
    categories: [
        "Character Design",
        "Environment Design",
        "Prop Design",
        "Creature Design",
        "Vehicle Design",
        "Keyframe Illustration",
        "Matte Painting",
        "UI/UX Design"
    ],
    // Keyboard shortcut: Ctrl + Shift + this key opens the admin modal
    adminShortcut: "ñ",
    // Backend API base URL (configurable for development / production)
    apiUrl: (typeof window !== "undefined" && (window.APP_CONFIG?.apiUrl || window.ENV?.API_URL)) || "http://localhost:8080/api",
    // Firebase Web configuration for Admin Authentication
    firebase: {
        apiKey: "AIzaSyDO8jZQ9cq3uJsRpNkbAIK-_x2owuqfD7w",
        authDomain: "dfcp-portfolio.firebaseapp.com",
        projectId: "dfcp-portfolio"
    }
};
