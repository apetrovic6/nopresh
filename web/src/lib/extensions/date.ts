import { timestampFromDate, type Timestamp } from "@bufbuild/protobuf/wkt";

declare global {
  interface Date {
    toTimestamp(): Timestamp | undefined
  }

}

Date.prototype.toTimestamp = function(this: Date) {
  return timestampFromDate(this);
}

export {};
