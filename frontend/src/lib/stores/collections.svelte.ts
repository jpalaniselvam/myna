
import { browser } from '$app/environment';

export class CollectionStore {
    collections = $state<string[]>([]);
    activeCollection = $state<string | null>(null);

    constructor() {
        if (browser) {
            const stored = localStorage.getItem('myna_collections');
            if (stored) {
                try {
                    this.collections = JSON.parse(stored);
                } catch (e) {
                    console.error("Failed to parse collections", e);
                    this.collections = [];
                }
            }

            const active = localStorage.getItem('myna_active_collection');
            if (active && this.collections.includes(active)) {
                this.activeCollection = active;
            }
        }
    }

    add(path: string) {
        if (!this.collections.includes(path)) {
            this.collections = [...this.collections, path];
            this.save();
        }
    }

    remove(path: string) {
        this.collections = this.collections.filter(c => c !== path);
        if (this.activeCollection === path) {
            this.activeCollection = null;
            if (browser) localStorage.removeItem('myna_active_collection');
        }
        this.save();
    }

    setActive(path: string) {
        if (this.collections.includes(path)) {
            this.activeCollection = path;
            if (browser) localStorage.setItem('myna_active_collection', path);
        }
    }

    save() {
        if (browser) {
            localStorage.setItem('myna_collections', JSON.stringify(this.collections));
        }
    }
}

export const collectionStore = new CollectionStore();
