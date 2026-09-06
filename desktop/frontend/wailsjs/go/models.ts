export namespace main {
	
	export class Finding {
	    file: string;
	    line: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Finding(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.line = source["line"];
	        this.message = source["message"];
	    }
	}

}

