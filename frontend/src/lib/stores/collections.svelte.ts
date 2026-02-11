


export class CollectionStore {
    collections = $state<string[]>([]);
    activeCollection = $state<string | null>(null);

    constructor() { }

    add(path: string) {
        if (!this.collections.includes(path)) {
            this.collections = [...this.collections, path];
        }
    }

    remove(path: string) {
        this.collections = this.collections.filter(c => c !== path);
        if (this.activeCollection === path) {
            this.activeCollection = null;
        }
    }

    setActive(path: string) {
        if (this.collections.includes(path)) {
            this.activeCollection = path;
        }
    }
}

export const collectionStore = new CollectionStore();
