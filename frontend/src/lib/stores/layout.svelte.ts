
export class LayoutStore {
    isSidebarOpen = $state(true);
    sidebarWidth = $state(256);

    constructor() { }

    toggleSidebar() {
        this.isSidebarOpen = !this.isSidebarOpen;
    }

    setSidebarOpen(value: boolean) {
        this.isSidebarOpen = value;
    }

    setSidebarWidth(width: number) {
        this.sidebarWidth = Math.max(200, Math.min(width, 600));
    }
}

export const layoutStore = new LayoutStore();
