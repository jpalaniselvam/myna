
export type TabKind = 'collection';

export interface Tab {
    id: string;
    title: string;
    kind: TabKind;
    metadata: any;
}

export class WorkspaceStore {
    tabs = $state<Tab[]>([]);
    activeTabId = $state<string | null>(null);

    get activeTab() {
        return this.tabs.find(t => t.id === this.activeTabId) || null;
    }

    openTab(kind: TabKind, id: string, title: string, metadata: any) {
        const existingTab = this.tabs.find(t => t.id === id);
        if (!existingTab) {
            this.tabs.push({ id, title, kind, metadata });
        }
        this.activeTabId = id;
    }

    closeTab(id: string) {
        const index = this.tabs.findIndex(t => t.id === id);
        if (index !== -1) {
            this.tabs.splice(index, 1);
            if (this.activeTabId === id) {
                if (this.tabs.length > 0) {
                    // Select the next tab, or the previous one if it was the last
                    const nextIndex = Math.min(index, this.tabs.length - 1);
                    this.activeTabId = this.tabs[nextIndex].id;
                } else {
                    this.activeTabId = null;
                }
            }
        }
    }

    setActiveTab(id: string) {
        if (this.tabs.find(t => t.id === id)) {
            this.activeTabId = id;
        }
    }
}

export const workspaceStore = new WorkspaceStore();
