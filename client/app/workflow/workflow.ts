export enum WorkflowType {
    IMPORT,
    EXPORT,
}

export class Workflow {
    public name: string
    public workflow: WorkflowType
    public refs: string
}
