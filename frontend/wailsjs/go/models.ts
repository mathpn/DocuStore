export namespace search {
	
	export class SearchResult {
	    DocID: string;
	    Title: string;
	    Identifier: string;
	    Type: string;
	    Score: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DocID = source["DocID"];
	        this.Title = source["Title"];
	        this.Identifier = source["Identifier"];
	        this.Type = source["Type"];
	        this.Score = source["Score"];
	    }
	}

}

