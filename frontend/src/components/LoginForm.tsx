import {
  Alert,
  Box,
  Button,
  Card,
  FormControl,
  FormLabel,
  Input,
  Typography,
} from "@mui/joy";
import { signInWithEmailAndPassword } from "firebase/auth";
import { useSetAtom } from "jotai";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { firebaseAuth } from "../api/firebase";
import { userAtom } from "../store/user";

interface FormElements extends HTMLFormControlsCollection {
  email: HTMLInputElement;
  password: HTMLInputElement;
}
interface SignInFormElement extends HTMLFormElement {
  readonly elements: FormElements;
}

export function LoginForm() {
  const navigate = useNavigate();
  const setUser = useSetAtom(userAtom);
  const [error, setError] = useState<string>("");
  function onSubmit(event: React.FormEvent<SignInFormElement>) {
    event.preventDefault();
    const formElements = event.currentTarget.elements;
    const data = {
      email: formElements.email.value,
      password: formElements.password.value,
    };
    signInWithEmailAndPassword(firebaseAuth, data.email, data.password)
      .then((userCredential) => {
        // setUser(userCredential.user);
        navigate("/");
      })
      .catch((error) => {
        const errorMessage = error.message ?? "unknown error";
        setError(errorMessage);
      });
  }
  return (
    <Box
      sx={{
        my: "auto",
        py: 2,
        display: "flex",
        flexDirection: "column",
        gap: 2,
        width: 400,
        maxWidth: "100%",
        mx: "auto",
        borderRadius: "sm",
      }}
    >
      <Card>
        <div>
          <Typography>Sign in to your account</Typography>
        </div>
        <form onSubmit={onSubmit} className="flex flex-col gap-2">
          <FormControl required>
            <FormLabel>Email</FormLabel>
            <Input type="email" name="email" />
          </FormControl>
          <FormControl required>
            <FormLabel>Password</FormLabel>
            <Input type="password" name="password" />
          </FormControl>
          {error.length > 0 && (
            <Alert color="danger" variant="soft">
              {error}
            </Alert>
          )}
          <Button type="submit" fullWidth>
            Sign in
          </Button>
        </form>
      </Card>
    </Box>
    // <Card title={"Sign in"}>
    //   <div>
    //     <label htmlFor="email" className="block mb-2 text-sm font-medium ">
    //       Your email
    //     </label>
    //     <input
    //       type="email"
    //       name="email"
    //       id="email"
    //       className="border border-gray-300 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
    //       placeholder="name@company.com"
    //       required
    //     />
    //   </div>
    //   <div>
    //     <label htmlFor="password" className="block mb-2 text-sm font-medium ">
    //       Your password
    //     </label>
    //     <input
    //       type="password"
    //       name="password"
    //       id="password"
    //       placeholder="••••••••"
    //       className=" border border-gray-300 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 "
    //       required
    //     />
    //   </div>
    //   <div className="flex items-start">
    //     <div className="flex items-start">
    //       <div className="flex items-center h-5">
    //         <input
    //           id="remember"
    //           type="checkbox"
    //           value=""
    //           className="w-4 h-4 border border-gray-300 rounded focus:ring-3 focus:ring-blue-300 "
    //           required
    //         />
    //       </div>
    //       <label htmlFor="remember" className="ml-2 text-sm font-medium ">
    //         Remember me
    //       </label>
    //     </div>
    //     <a href="#" className="ml-auto text-sm text-blue-700 hover:underline ">
    //       Lost Password?
    //     </a>
    //   </div>
    //   <Button type="submit" width="full">
    //     Login to your account
    //   </Button>
    //   <div className="text-sm font-medium text-gray-500 ">
    //     Not registered?{" "}
    //     <a href="#" className="text-blue-700 hover:underline">
    //       Create account
    //     </a>
    //   </div>
    // </Card>
  );
}
