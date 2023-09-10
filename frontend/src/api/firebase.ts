import { FirebaseOptions, initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
const firebaseApiKey = import.meta.env.VITE_FIREBASE_API_KEY;
const firebaseConfig: FirebaseOptions = {
  apiKey: firebaseApiKey,
};

export const firebaseApp = initializeApp(firebaseConfig);
export const firebaseAuth = getAuth(firebaseApp);
