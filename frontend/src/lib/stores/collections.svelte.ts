


import type { collection } from '../../../wailsjs/go/models';

export class CollectionStore {
    collections = $state<string[]>([]);
    collectionDetails = $state<Record<string, collection.CollectionResponse>>({});
    activeCollection = $state<string | null>(null);

    constructor() { }

    private normalize(path: string): string {
        return path.replace(/\\/g, '/');
    }

    add(path: string) {
        const normalized = this.normalize(path);
        if (!this.collections.some(c => this.normalize(c) === normalized)) {
            this.collections = [...this.collections, path];
        }
    }

    remove(path: string) {
        const normalized = this.normalize(path);
        this.collections = this.collections.filter(c => this.normalize(c) !== normalized);
        delete this.collectionDetails[normalized];
        if (this.activeCollection && this.normalize(this.activeCollection) === normalized) {
            this.activeCollection = null;
        }
    }

    setDetails(path: string, details: collection.CollectionResponse) {
        const normalized = this.normalize(path);
        this.collectionDetails[normalized] = details;
    }

    setActive(path: string | null) {
        if (!path) {
            this.activeCollection = null;
            return;
        }
        const normalized = this.normalize(path);
        const match = this.collections.find(c => this.normalize(c) === normalized);
        if (match) {
            this.activeCollection = match;
        }
    }
}


export const collectionStore = new CollectionStore();
