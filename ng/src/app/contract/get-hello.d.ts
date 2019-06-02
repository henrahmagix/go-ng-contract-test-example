export interface HelloResponse {
  method: string;
  path: string;
  response: Response;
}
export interface Response {
  msg: string;
  obj: Obj;
  num: number;
  arr?: (string)[] | null;
}
export interface Obj {
}
