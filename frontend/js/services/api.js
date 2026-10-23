import { config } from '../config.js';

const API_URL = (config.apiUrl || "http://localhost:8080/api").replace(/\/+$/, '');

export async function fetchArtworks() {
    const res = await fetch(`${API_URL}/artworks`);
    if (!res.ok) throw new Error('Failed to fetch artworks');
    return res.json();
}

export async function createArtwork(formData, token) {
    const res = await fetch(`${API_URL}/artworks`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });
    if (!res.ok) throw new Error('Failed to create artwork');
    return res.json();
}

export async function updateArtwork(id, data, token) {
    const res = await fetch(`${API_URL}/artworks/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(data)
    });
    if (!res.ok) throw new Error('Failed to update artwork');
    return res.json();
}

export async function deleteArtwork(id, token) {
    const res = await fetch(`${API_URL}/artworks/${id}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    if (!res.ok) throw new Error('Failed to delete artwork');
    return res.json();
}

export async function fetchResumeURL() {
    const res = await fetch(`${API_URL}/resume`);
    if (!res.ok) throw new Error('Failed to fetch resume URL');
    return res.json();
}

export async function uploadResume(formData, token) {
    const res = await fetch(`${API_URL}/resume`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    });
    if (!res.ok) throw new Error('Failed to upload resume');
    return res.json();
}
