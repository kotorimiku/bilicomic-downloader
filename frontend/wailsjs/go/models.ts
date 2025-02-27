export namespace bilicomicdownloader {
	
	export class BookInfo {
	    Author: string[];
	    Description: string;
	    Genre: string[];
	    Title: string;
	    Cover: string;
	
	    static createFrom(source: any = {}) {
	        return new BookInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Author = source["Author"];
	        this.Description = source["Description"];
	        this.Genre = source["Genre"];
	        this.Title = source["Title"];
	        this.Cover = source["Cover"];
	    }
	}
	export class Chapter {
	    Title: string;
	    Url: string;
	
	    static createFrom(source: any = {}) {
	        return new Chapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Url = source["Url"];
	    }
	}
	export class Config {
	    urlBase: string;
	    outputPath: string;
	    packageType: string;
	    imageFormat: string;
	    namingStyle: string;
	    cookie: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.urlBase = source["urlBase"];
	        this.outputPath = source["outputPath"];
	        this.packageType = source["packageType"];
	        this.imageFormat = source["imageFormat"];
	        this.namingStyle = source["namingStyle"];
	        this.cookie = source["cookie"];
	    }
	}
	export class Volume {
	    Title: string;
	    Cover: string;
	    Chapters: Chapter[];
	
	    static createFrom(source: any = {}) {
	        return new Volume(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Cover = source["Cover"];
	        this.Chapters = this.convertValues(source["Chapters"], Chapter);
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
	export class DownloaderSingle {
	    BookInfo?: BookInfo;
	    Volume?: Volume;
	    Index: number;
	    Progress: number;
	    Fail: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DownloaderSingle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BookInfo = this.convertValues(source["BookInfo"], BookInfo);
	        this.Volume = this.convertValues(source["Volume"], Volume);
	        this.Index = source["Index"];
	        this.Progress = source["Progress"];
	        this.Fail = source["Fail"];
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

}

