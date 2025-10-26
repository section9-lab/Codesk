export namespace checkpoint {
	
	export class CheckpointMetadata {
	    totalTokens: number;
	    modelUsed: string;
	    userPrompt: string;
	    fileChanges: number;
	    snapshotSize: number;
	
	    static createFrom(source: any = {}) {
	        return new CheckpointMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalTokens = source["totalTokens"];
	        this.modelUsed = source["modelUsed"];
	        this.userPrompt = source["userPrompt"];
	        this.fileChanges = source["fileChanges"];
	        this.snapshotSize = source["snapshotSize"];
	    }
	}
	export class Checkpoint {
	    id: string;
	    sessionId: string;
	    projectId: string;
	    messageIndex: number;
	    // Go type: time
	    timestamp: any;
	    description?: string;
	    parentCheckpointId?: string;
	    metadata: CheckpointMetadata;
	
	    static createFrom(source: any = {}) {
	        return new Checkpoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sessionId = source["sessionId"];
	        this.projectId = source["projectId"];
	        this.messageIndex = source["messageIndex"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.description = source["description"];
	        this.parentCheckpointId = source["parentCheckpointId"];
	        this.metadata = this.convertValues(source["metadata"], CheckpointMetadata);
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
	export class FileDiff {
	    path: string;
	    additions: number;
	    deletions: number;
	    diffContent?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDiff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.additions = source["additions"];
	        this.deletions = source["deletions"];
	        this.diffContent = source["diffContent"];
	    }
	}
	export class CheckpointDiff {
	    fromCheckpointId: string;
	    toCheckpointId: string;
	    modifiedFiles: FileDiff[];
	    addedFiles: string[];
	    deletedFiles: string[];
	    tokenDelta: number;
	
	    static createFrom(source: any = {}) {
	        return new CheckpointDiff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fromCheckpointId = source["fromCheckpointId"];
	        this.toCheckpointId = source["toCheckpointId"];
	        this.modifiedFiles = this.convertValues(source["modifiedFiles"], FileDiff);
	        this.addedFiles = source["addedFiles"];
	        this.deletedFiles = source["deletedFiles"];
	        this.tokenDelta = source["tokenDelta"];
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
	
	export class CheckpointResult {
	    checkpoint: Checkpoint;
	    filesProcessed: number;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new CheckpointResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.checkpoint = this.convertValues(source["checkpoint"], Checkpoint);
	        this.filesProcessed = source["filesProcessed"];
	        this.warnings = source["warnings"];
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
	export class CheckpointSettings {
	    autoCheckpointEnabled: boolean;
	    checkpointStrategy: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckpointSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.autoCheckpointEnabled = source["autoCheckpointEnabled"];
	        this.checkpointStrategy = source["checkpointStrategy"];
	    }
	}
	
	export class TimelineNode {
	    checkpoint: Checkpoint;
	    children: TimelineNode[];
	    fileSnapshotIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new TimelineNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.checkpoint = this.convertValues(source["checkpoint"], Checkpoint);
	        this.children = this.convertValues(source["children"], TimelineNode);
	        this.fileSnapshotIds = source["fileSnapshotIds"];
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
	export class SessionTimeline {
	    sessionId: string;
	    rootNode?: TimelineNode;
	    currentCheckpointId?: string;
	    autoCheckpointEnabled: boolean;
	    checkpointStrategy: string;
	    totalCheckpoints: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionTimeline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.rootNode = this.convertValues(source["rootNode"], TimelineNode);
	        this.currentCheckpointId = source["currentCheckpointId"];
	        this.autoCheckpointEnabled = source["autoCheckpointEnabled"];
	        this.checkpointStrategy = source["checkpointStrategy"];
	        this.totalCheckpoints = source["totalCheckpoints"];
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

export namespace claude {
	
	export class ClaudeMdFile {
	    relative_path: string;
	    absolute_path: string;
	    size: number;
	    modified: number;
	
	    static createFrom(source: any = {}) {
	        return new ClaudeMdFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.relative_path = source["relative_path"];
	        this.absolute_path = source["absolute_path"];
	        this.size = source["size"];
	        this.modified = source["modified"];
	    }
	}
	export class ExecuteResult {
	    session_id: string;
	    project_path: string;
	    pid: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ExecuteResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.session_id = source["session_id"];
	        this.project_path = source["project_path"];
	        this.pid = source["pid"];
	        this.message = source["message"];
	    }
	}
	export class FileEntry {
	    name: string;
	    path: string;
	    is_directory: boolean;
	    size: number;
	    extension?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.is_directory = source["is_directory"];
	        this.size = source["size"];
	        this.extension = source["extension"];
	    }
	}

}

export namespace model {
	
	export class Agent {
	    id?: number;
	    name: string;
	    icon: string;
	    system_prompt: string;
	    default_task?: string;
	    model: string;
	    enable_file_read: boolean;
	    enable_file_write: boolean;
	    enable_network: boolean;
	    hooks?: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Agent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.system_prompt = source["system_prompt"];
	        this.default_task = source["default_task"];
	        this.model = source["model"];
	        this.enable_file_read = source["enable_file_read"];
	        this.enable_file_write = source["enable_file_write"];
	        this.enable_network = source["enable_network"];
	        this.hooks = source["hooks"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class AgentData {
	    name: string;
	    icon: string;
	    system_prompt: string;
	    default_task?: string;
	    model: string;
	    hooks?: string;
	
	    static createFrom(source: any = {}) {
	        return new AgentData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.system_prompt = source["system_prompt"];
	        this.default_task = source["default_task"];
	        this.model = source["model"];
	        this.hooks = source["hooks"];
	    }
	}
	export class AgentExport {
	    version: number;
	    // Go type: time
	    exported_at: any;
	    agent: AgentData;
	
	    static createFrom(source: any = {}) {
	        return new AgentExport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.exported_at = this.convertValues(source["exported_at"], null);
	        this.agent = this.convertValues(source["agent"], AgentData);
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
	export class AgentRun {
	    id?: number;
	    agent_id: number;
	    agent_name: string;
	    agent_icon: string;
	    task: string;
	    model: string;
	    project_path: string;
	    session_id: string;
	    status: string;
	    pid?: number;
	    // Go type: time
	    process_started_at?: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    completed_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new AgentRun(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.agent_id = source["agent_id"];
	        this.agent_name = source["agent_name"];
	        this.agent_icon = source["agent_icon"];
	        this.task = source["task"];
	        this.model = source["model"];
	        this.project_path = source["project_path"];
	        this.session_id = source["session_id"];
	        this.status = source["status"];
	        this.pid = source["pid"];
	        this.process_started_at = this.convertValues(source["process_started_at"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.completed_at = this.convertValues(source["completed_at"], null);
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
	export class ClaudeSettings {
	    data: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ClaudeSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	    }
	}
	export class DailyUsageStats {
	    date: string;
	    sessions: number;
	    tokens: number;
	    cost_usd: number;
	
	    static createFrom(source: any = {}) {
	        return new DailyUsageStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.sessions = source["sessions"];
	        this.tokens = source["tokens"];
	        this.cost_usd = source["cost_usd"];
	    }
	}
	export class MCPServer {
	    name: string;
	    command: string;
	    args: string[];
	    env: Record<string, string>;
	    disabled: boolean;
	    auto_approve: string[];
	
	    static createFrom(source: any = {}) {
	        return new MCPServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.env = source["env"];
	        this.disabled = source["disabled"];
	        this.auto_approve = source["auto_approve"];
	    }
	}
	export class MCPProjectConfig {
	    project_path: string;
	    servers: Record<string, MCPServer>;
	    choices: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new MCPProjectConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_path = source["project_path"];
	        this.servers = this.convertValues(source["servers"], MCPServer, true);
	        this.choices = source["choices"];
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
	
	export class MCPServerStatus {
	    name: string;
	    status: string;
	    pid?: number;
	    error: string;
	    connected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MCPServerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.pid = source["pid"];
	        this.error = source["error"];
	        this.connected = source["connected"];
	    }
	}
	export class MessageContent {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new MessageContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class Project {
	    id: string;
	    path: string;
	    sessions: string[];
	    created_at: number;
	    most_recent_session?: number;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.sessions = source["sessions"];
	        this.created_at = source["created_at"];
	        this.most_recent_session = source["most_recent_session"];
	    }
	}
	export class ProxySettings {
	    enabled: boolean;
	    http_proxy?: string;
	    https_proxy?: string;
	    no_proxy?: string;
	    all_proxy?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxySettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.http_proxy = source["http_proxy"];
	        this.https_proxy = source["https_proxy"];
	        this.no_proxy = source["no_proxy"];
	        this.all_proxy = source["all_proxy"];
	    }
	}
	export class SQLResult {
	    rows_affected: number;
	    last_insert_id: number;
	    rows?: any[];
	
	    static createFrom(source: any = {}) {
	        return new SQLResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rows_affected = source["rows_affected"];
	        this.last_insert_id = source["last_insert_id"];
	        this.rows = source["rows"];
	    }
	}
	export class Session {
	    id: string;
	    project_id: string;
	    project_path: string;
	    todo_data: any;
	    created_at: number;
	    first_message?: string;
	    message_timestamp?: string;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.project_path = source["project_path"];
	        this.todo_data = source["todo_data"];
	        this.created_at = source["created_at"];
	        this.first_message = source["first_message"];
	        this.message_timestamp = source["message_timestamp"];
	    }
	}
	export class SessionMessage {
	    type: string;
	    message?: MessageContent;
	    timestamp?: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.message = this.convertValues(source["message"], MessageContent);
	        this.timestamp = source["timestamp"];
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
	export class SessionHistory {
	    session_id: string;
	    messages: SessionMessage[];
	
	    static createFrom(source: any = {}) {
	        return new SessionHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.session_id = source["session_id"];
	        this.messages = this.convertValues(source["messages"], SessionMessage);
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
	
	export class SessionStats {
	    session_id: string;
	    message_count: number;
	    total_tokens: number;
	    input_tokens: number;
	    output_tokens: number;
	    cost_usd: number;
	    duration: number;
	    // Go type: time
	    start_time: any;
	    // Go type: time
	    end_time?: any;
	
	    static createFrom(source: any = {}) {
	        return new SessionStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.session_id = source["session_id"];
	        this.message_count = source["message_count"];
	        this.total_tokens = source["total_tokens"];
	        this.input_tokens = source["input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.cost_usd = source["cost_usd"];
	        this.duration = source["duration"];
	        this.start_time = this.convertValues(source["start_time"], null);
	        this.end_time = this.convertValues(source["end_time"], null);
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
	export class SlashCommand {
	    id?: number;
	    name: string;
	    description: string;
	    command: string;
	    icon: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new SlashCommand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.command = source["command"];
	        this.icon = source["icon"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class TableInfo {
	    name: string;
	    row_count: number;
	    columns: string[];
	
	    static createFrom(source: any = {}) {
	        return new TableInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.row_count = source["row_count"];
	        this.columns = source["columns"];
	    }
	}
	export class UsageByDateRange {
	    start_date: string;
	    end_date: string;
	    total_sessions: number;
	    total_tokens: number;
	    total_cost_usd: number;
	    daily_stats: DailyUsageStats[];
	
	    static createFrom(source: any = {}) {
	        return new UsageByDateRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.total_sessions = source["total_sessions"];
	        this.total_tokens = source["total_tokens"];
	        this.total_cost_usd = source["total_cost_usd"];
	        this.daily_stats = this.convertValues(source["daily_stats"], DailyUsageStats);
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
	export class UsageStats {
	    total_sessions: number;
	    total_messages: number;
	    total_tokens: number;
	    total_cost_usd: number;
	    average_tokens: number;
	    average_cost_usd: number;
	
	    static createFrom(source: any = {}) {
	        return new UsageStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_sessions = source["total_sessions"];
	        this.total_messages = source["total_messages"];
	        this.total_tokens = source["total_tokens"];
	        this.total_cost_usd = source["total_cost_usd"];
	        this.average_tokens = source["average_tokens"];
	        this.average_cost_usd = source["average_cost_usd"];
	    }
	}

}

