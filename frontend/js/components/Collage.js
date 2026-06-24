export class Collage {
    /**
     * @param {HTMLElement} container - The DOM element where the collage will be injected
     * @param {Array} artworks - Array of artworks (ARTWORK model)
     */
    constructor(container, artworks) {
        this.container = container;
        this.artworks = artworks;
    }

    render() {
        this.container.innerHTML = ''; // Clear the container

        this.artworks.forEach((art) => {
            const item = document.createElement('div');
            item.className = 'collage__item';

            const img = document.createElement('img');
            img.src = art.image_url;
            img.alt = art.title;
            img.className = 'collage__image';
            img.loading = 'lazy'; // Native lazy loading

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
        });
    }
}
