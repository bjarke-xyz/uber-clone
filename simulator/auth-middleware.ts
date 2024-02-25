import { NextFunction, Request, Response } from "express";
type AuthRole = "ADMIN";
export function authMiddleware(role: AuthRole) {
  return (req: Request, res: Response, next: NextFunction) => {
    const authorizationHeader = req.headers["authorization"];
    if (!authorizationHeader) {
      res.status(401);
      return res.send("missing authorization header");
    }
    const token = authorizationHeader.split("Bearer ")[1];
    if (!token) {
      res.status(401);
      return res.send("bad bearer token format");
    }
    // getAuthClient().validateToken(
    //   {
    //     token: token,
    //     audience: "",
    //   },
    //   (err, resp) => {
    //     if (err) {
    //       res.status(401);
    //       return res.json(err);
    //     } else {
    //       if (resp.role?.toUpperCase() !== role.toUpperCase()) {
    //         res.status(403);
    //         return res.send("invalid role");
    //       }
    //       setAuthToken(req, resp);
    //       next();
    //     }
    //   }
    // );
  };
}
