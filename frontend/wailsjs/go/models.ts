export namespace collection {
	
	export class CollectionResponse {
	    description: string;
	    environments: string[];
	    pre: Record<string, any>;
	    actions: Record<string, any>;
	    credentials: types.Credentials;
	
	    static createFrom(source: any = {}) {
	        return new CollectionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.environments = source["environments"];
	        this.pre = source["pre"];
	        this.actions = source["actions"];
	        this.credentials = this.convertValues(source["credentials"], types.Credentials);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateActionInput {
	    collection_path: string;
	    sub_path: string;
	    file_name: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new CreateActionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.collection_path = source["collection_path"];
	        this.sub_path = source["sub_path"];
	        this.file_name = source["file_name"];
	        this.data = source["data"];
	    }
	}
	export class GetActionInput {
	    collection_path: string;
	    sub_path: string;
	    file_name: string;
	
	    static createFrom(source: any = {}) {
	        return new GetActionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.collection_path = source["collection_path"];
	        this.sub_path = source["sub_path"];
	        this.file_name = source["file_name"];
	    }
	}
	export class UpdateActionInput {
	    collection_path: string;
	    sub_path: string;
	    file_name: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new UpdateActionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.collection_path = source["collection_path"];
	        this.sub_path = source["sub_path"];
	        this.file_name = source["file_name"];
	        this.data = source["data"];
	    }
	}

}

export namespace types {
	
	export class Credentials {
	    region?: string;
	    profile?: string;
	    role_arn?: string;
	
	    static createFrom(source: any = {}) {
	        return new Credentials(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.region = source["region"];
	        this.profile = source["profile"];
	        this.role_arn = source["role_arn"];
	    }
	}

}

