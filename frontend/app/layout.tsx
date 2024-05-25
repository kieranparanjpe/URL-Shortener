import type { Metadata } from "next";
import "./reset.css"
import "./globals.css";

export const metadata: Metadata = {
  title: "URL Shortener",
  description: "Shorten your urls",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
