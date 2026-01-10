const Auth = {
    isAuthenticated() {
        return sessionStorage.getItem('siem_credentials') !== null;
    },
    
    async login(username, password) {
        try {
            return await API.testAuth(username, password);
        } catch (error) {
            console.error('Login error:', error);
            return false;
        }
    },
    
    logout() {
        sessionStorage.removeItem('siem_credentials');
        window.location.href = 'login.html';
    },
    
    requireAuth() {
        if (!this.isAuthenticated()) {
            window.location.href = 'login.html';
            return false;
        }
        return true;
    },

    getCredentials() {
        const stored = sessionStorage.getItem('siem_credentials');
        if (!stored) return null;
        try {
            const decoded = atob(stored);
            const [username, password] = decoded.split(':');
            return { username, password };
        } catch (e) {
            console.error('Invalid credentials format', e);
            return null;
        }
    }
};

document.addEventListener('DOMContentLoaded', () => {
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => Auth.logout());
    }
});

window.Auth = Auth;
