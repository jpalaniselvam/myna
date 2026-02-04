export namespace collection {
	
	export class CollectionResponse {
	    description: string;
	    environments: string[];
	    pre: Record<string, any>;
	    actions: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new CollectionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.description = source["description"];
	        this.environments = source["environments"];
	        this.pre = source["pre"];
	        this.actions = source["actions"];
	    }
	}
	export class CreateActionInput {
	    collection_path: string;
	    sub_path: string;
	    file_name: string;
	    data: number[];
	
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
	    data: number[];
	
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

