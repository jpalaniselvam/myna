
export class LayoutStore {
    isSidebarOpen = $state(true);

    constructor() { }

    toggleSidebar() {
        this.isSidebarOpen = !this.isSidebarOpen;
    }

    setSidebarOpen(value: boolean) {
        this.isSidebarOpen = value;
    }
}

export const layoutStore = new LayoutStore();
