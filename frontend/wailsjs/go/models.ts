export namespace config {
	
	export class ClientConnection {
	    connected: boolean;
	    permission: string;
	    connectedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connected = source["connected"];
	        this.permission = source["permission"];
	        this.connectedAt = source["connectedAt"];
	    }
	}
	export class AppConfig {
	    version: number;
	    clients: Record<string, ClientConnection>;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.clients = this.convertValues(source["clients"], ClientConnection, true);
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

export namespace discovery {
	
	export class ClientInfo {
	    id: string;
	    name: string;
	    detected: boolean;
	    connected: boolean;
	    permission: string;
	    executable: string;
	    configPaths: string[];
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.detected = source["detected"];
	        this.connected = source["connected"];
	        this.permission = source["permission"];
	        this.executable = source["executable"];
	        this.configPaths = source["configPaths"];
	        this.status = source["status"];
	    }
	}
	export class MCPServer {
	    name: string;
	    clientID: string;
	    client: string;
	    transport: string;
	    command: string;
	    args: string;
	    url: string;
	    env: string;
	    cwd: string;
	    status: string;
	    auth: string;
	    sourcePath: string;
	    origin: string;
	    originPath: string;
	
	    static createFrom(source: any = {}) {
	        return new MCPServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.clientID = source["clientID"];
	        this.client = source["client"];
	        this.transport = source["transport"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.url = source["url"];
	        this.env = source["env"];
	        this.cwd = source["cwd"];
	        this.status = source["status"];
	        this.auth = source["auth"];
	        this.sourcePath = source["sourcePath"];
	        this.origin = source["origin"];
	        this.originPath = source["originPath"];
	    }
	}
	export class SkillAsset {
	    name: string;
	    clientID: string;
	    client: string;
	    description: string;
	    sourcePath: string;
	    scope: string;
	
	    static createFrom(source: any = {}) {
	        return new SkillAsset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.clientID = source["clientID"];
	        this.client = source["client"];
	        this.description = source["description"];
	        this.sourcePath = source["sourcePath"];
	        this.scope = source["scope"];
	    }
	}

}

