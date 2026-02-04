
export class EnvironmentStore {
    activeEnvironment = $state<string | null>(null);

    setActive(env: string | null) {
        this.activeEnvironment = env;
    }
}

export const environmentStore = new EnvironmentStore();
