export enum MappingType {
    IMPORT,
    EXPORT,
}

export class Mapping {
    public name: string
    public mapping: MappingType
    public refs: string
}
