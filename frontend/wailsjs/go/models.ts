export namespace collection {
	
	export class CollectionResponse {
	    environments: string[];
	    pre: Record<string, any>;
	    actions: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new CollectionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.environments = source["environments"];
	        this.pre = source["pre"];
	        this.actions = source["actions"];
	    }
	}

}

