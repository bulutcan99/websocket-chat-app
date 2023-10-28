import { AllPages } from "@/utils/types/enums";
import React, { ReactNode } from "react";

interface Props {
  children: ReactNode;
  page?: string;
}

const Layout = ({ children, page }: Props) => {
  if (page === AllPages.LOGIN || page === AllPages.REGISTER) {
    return <main>{children}</main>;
  }

  return (
    <>
      <header>Header</header>
      <main>{children}</main>
      <footer>Footer</footer>
    </>
  );
};

export default Layout;
