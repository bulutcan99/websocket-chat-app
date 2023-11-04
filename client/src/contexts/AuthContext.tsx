import React, {
  useState,
  useContext,
  createContext,
  ReactNode,
  Dispatch,
  SetStateAction,
  FunctionComponent,
} from "react";

interface AuthContextProps {
  isAuthenticated?: boolean;
  setIsAuthenticated?: Dispatch<SetStateAction<boolean>>;
  login?: () => void;
  logout?: () => void;
}

interface AuthProviderProps {
  children: ReactNode;
}

const AuthContext = createContext<AuthContextProps>({});

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const login = () => {
    setIsAuthenticated(true);
  };

  const logout = () => {
    setIsAuthenticated(false);
  };

  const authContextValue: AuthContextProps = {
    isAuthenticated,
    setIsAuthenticated,
    login,
    logout,
  };

  return (
    <AuthContext.Provider value={authContextValue}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuthContext = () => useContext(AuthContext);
