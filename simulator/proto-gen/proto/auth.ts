/* eslint-disable */
import {
  CallOptions,
  ChannelCredentials,
  Client,
  ClientOptions,
  ClientUnaryCall,
  handleUnaryCall,
  makeGenericClientConstructor,
  Metadata,
  ServiceError,
  UntypedServiceImplementation,
} from "@grpc/grpc-js";
import _m0 from "protobufjs/minimal";
import { Any } from "../google/protobuf/any";

export const protobufPackage = "auth";

export interface AuthToken {
  authTime: number;
  issuer: string;
  audience: string;
  expires: number;
  issuedAt: number;
  subject: string;
  UID: string;
  claims: { [key: string]: Any };
  products: string[];
  role: string;
  groups: string[];
}

export interface AuthToken_ClaimsEntry {
  key: string;
  value: Any | undefined;
}

export interface ValidateTokenRequest {
  token: string;
  audience: string;
}

function createBaseAuthToken(): AuthToken {
  return {
    authTime: 0,
    issuer: "",
    audience: "",
    expires: 0,
    issuedAt: 0,
    subject: "",
    UID: "",
    claims: {},
    products: [],
    role: "",
    groups: [],
  };
}

export const AuthToken = {
  encode(message: AuthToken, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authTime !== 0) {
      writer.uint32(9).double(message.authTime);
    }
    if (message.issuer !== "") {
      writer.uint32(18).string(message.issuer);
    }
    if (message.audience !== "") {
      writer.uint32(26).string(message.audience);
    }
    if (message.expires !== 0) {
      writer.uint32(33).double(message.expires);
    }
    if (message.issuedAt !== 0) {
      writer.uint32(41).double(message.issuedAt);
    }
    if (message.subject !== "") {
      writer.uint32(50).string(message.subject);
    }
    if (message.UID !== "") {
      writer.uint32(58).string(message.UID);
    }
    Object.entries(message.claims).forEach(([key, value]) => {
      AuthToken_ClaimsEntry.encode({ key: key as any, value }, writer.uint32(66).fork()).ldelim();
    });
    for (const v of message.products) {
      writer.uint32(74).string(v!);
    }
    if (message.role !== "") {
      writer.uint32(82).string(message.role);
    }
    for (const v of message.groups) {
      writer.uint32(90).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthToken {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.authTime = reader.double();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.issuer = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.audience = reader.string();
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.expires = reader.double();
          continue;
        case 5:
          if (tag !== 41) {
            break;
          }

          message.issuedAt = reader.double();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.subject = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.UID = reader.string();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          const entry8 = AuthToken_ClaimsEntry.decode(reader, reader.uint32());
          if (entry8.value !== undefined) {
            message.claims[entry8.key] = entry8.value;
          }
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.products.push(reader.string());
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.role = reader.string();
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.groups.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthToken {
    return {
      authTime: isSet(object.authTime) ? Number(object.authTime) : 0,
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      audience: isSet(object.audience) ? String(object.audience) : "",
      expires: isSet(object.expires) ? Number(object.expires) : 0,
      issuedAt: isSet(object.issuedAt) ? Number(object.issuedAt) : 0,
      subject: isSet(object.subject) ? String(object.subject) : "",
      UID: isSet(object.UID) ? String(object.UID) : "",
      claims: isObject(object.claims)
        ? Object.entries(object.claims).reduce<{ [key: string]: Any }>((acc, [key, value]) => {
          acc[key] = Any.fromJSON(value);
          return acc;
        }, {})
        : {},
      products: Array.isArray(object?.products) ? object.products.map((e: any) => String(e)) : [],
      role: isSet(object.role) ? String(object.role) : "",
      groups: Array.isArray(object?.groups) ? object.groups.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: AuthToken): unknown {
    const obj: any = {};
    if (message.authTime !== 0) {
      obj.authTime = message.authTime;
    }
    if (message.issuer !== "") {
      obj.issuer = message.issuer;
    }
    if (message.audience !== "") {
      obj.audience = message.audience;
    }
    if (message.expires !== 0) {
      obj.expires = message.expires;
    }
    if (message.issuedAt !== 0) {
      obj.issuedAt = message.issuedAt;
    }
    if (message.subject !== "") {
      obj.subject = message.subject;
    }
    if (message.UID !== "") {
      obj.UID = message.UID;
    }
    if (message.claims) {
      const entries = Object.entries(message.claims);
      if (entries.length > 0) {
        obj.claims = {};
        entries.forEach(([k, v]) => {
          obj.claims[k] = Any.toJSON(v);
        });
      }
    }
    if (message.products?.length) {
      obj.products = message.products;
    }
    if (message.role !== "") {
      obj.role = message.role;
    }
    if (message.groups?.length) {
      obj.groups = message.groups;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AuthToken>, I>>(base?: I): AuthToken {
    return AuthToken.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AuthToken>, I>>(object: I): AuthToken {
    const message = createBaseAuthToken();
    message.authTime = object.authTime ?? 0;
    message.issuer = object.issuer ?? "";
    message.audience = object.audience ?? "";
    message.expires = object.expires ?? 0;
    message.issuedAt = object.issuedAt ?? 0;
    message.subject = object.subject ?? "";
    message.UID = object.UID ?? "";
    message.claims = Object.entries(object.claims ?? {}).reduce<{ [key: string]: Any }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = Any.fromPartial(value);
      }
      return acc;
    }, {});
    message.products = object.products?.map((e) => e) || [];
    message.role = object.role ?? "";
    message.groups = object.groups?.map((e) => e) || [];
    return message;
  },
};

function createBaseAuthToken_ClaimsEntry(): AuthToken_ClaimsEntry {
  return { key: "", value: undefined };
}

export const AuthToken_ClaimsEntry = {
  encode(message: AuthToken_ClaimsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Any.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthToken_ClaimsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthToken_ClaimsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = Any.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthToken_ClaimsEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Any.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: AuthToken_ClaimsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = Any.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AuthToken_ClaimsEntry>, I>>(base?: I): AuthToken_ClaimsEntry {
    return AuthToken_ClaimsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AuthToken_ClaimsEntry>, I>>(object: I): AuthToken_ClaimsEntry {
    const message = createBaseAuthToken_ClaimsEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null) ? Any.fromPartial(object.value) : undefined;
    return message;
  },
};

function createBaseValidateTokenRequest(): ValidateTokenRequest {
  return { token: "", audience: "" };
}

export const ValidateTokenRequest = {
  encode(message: ValidateTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.audience !== "") {
      writer.uint32(18).string(message.audience);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidateTokenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidateTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.token = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.audience = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ValidateTokenRequest {
    return {
      token: isSet(object.token) ? String(object.token) : "",
      audience: isSet(object.audience) ? String(object.audience) : "",
    };
  },

  toJSON(message: ValidateTokenRequest): unknown {
    const obj: any = {};
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.audience !== "") {
      obj.audience = message.audience;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ValidateTokenRequest>, I>>(base?: I): ValidateTokenRequest {
    return ValidateTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ValidateTokenRequest>, I>>(object: I): ValidateTokenRequest {
    const message = createBaseValidateTokenRequest();
    message.token = object.token ?? "";
    message.audience = object.audience ?? "";
    return message;
  },
};

export type AuthService = typeof AuthService;
export const AuthService = {
  validateToken: {
    path: "/auth.Auth/ValidateToken",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: ValidateTokenRequest) => Buffer.from(ValidateTokenRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => ValidateTokenRequest.decode(value),
    responseSerialize: (value: AuthToken) => Buffer.from(AuthToken.encode(value).finish()),
    responseDeserialize: (value: Buffer) => AuthToken.decode(value),
  },
} as const;

export interface AuthServer extends UntypedServiceImplementation {
  validateToken: handleUnaryCall<ValidateTokenRequest, AuthToken>;
}

export interface AuthClient extends Client {
  validateToken(
    request: ValidateTokenRequest,
    callback: (error: ServiceError | null, response: AuthToken) => void,
  ): ClientUnaryCall;
  validateToken(
    request: ValidateTokenRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: AuthToken) => void,
  ): ClientUnaryCall;
  validateToken(
    request: ValidateTokenRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: AuthToken) => void,
  ): ClientUnaryCall;
}

export const AuthClient = makeGenericClientConstructor(AuthService, "auth.Auth") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): AuthClient;
  service: typeof AuthService;
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
