import AuthState from "./Authentication/AuthState";
import API from "./util/API";

export default async function Home() {
  API.init();
  
  return (
    <main>
        <div>
            <div className={"content-start text-center justify-start bg-white shadow-xl shadow-gray-400"}
                 style={{marginLeft: "8%", marginRight: "8%", marginTop: "2%", minHeight: "90vh", borderRadius: "20px"}}>
                <div className="flex items-center justify-items-center w-full h-32">
                    <h2 className="text-black text-left font-bold text-7xl" style={{paddingInline: "1vw", flex: "1"}}>URL Shortener</h2>
                </div>
                <div className="flex items-center justify-items-center w-full h-2 bg-black"></div>
                <AuthState/>
            </div>
        </div>
    </main>
  );
}
