import { useState } from "react";

function App() {
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  const handleClick = async () => {
    setLoading(true);
    setError(null);
    setResults(null);

    try {
      const [tab] = await chrome.tabs.query({
        active: true,
        currentWindow: true,
      });

      if (tab.id) {
        const [injectionResult] = await chrome.scripting.executeScript({
          target: { tabId: tab.id },
          func: () => {
            const video = document.querySelector("video");
            if (!video) {
              return { error: "No video found on the page" };
            }
            const canvas = document.createElement("canvas");
            canvas.width = video.videoWidth;
            canvas.height = video.videoHeight;
            const ctx = canvas.getContext("2d");
            if (!ctx) {
              return { error: "Could not get canvas context" };
            }
            ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
            return { dataUrl: canvas.toDataURL("image/png") };
          },
        });

        const result = injectionResult?.result as {
          dataUrl?: string;
          error?: string;
        };

        if (result?.error) {
          setError(result.error);
          setLoading(false);
          return;
        }

        if (!result?.dataUrl) {
          setError("Could not get image data");
          setLoading(false);
          return;
        }

        const { dataUrl } = result;

        const blob = await fetch(dataUrl).then((res) => res.blob());

        const formData = new FormData();
        formData.append("image", blob, "screenshot.png");

        const apiResponse = await fetch(
          "https://murmuring-tundra-32346-4f714b7b806b.herokuapp.com/identify-clothing",
          {
            method: "POST",
            body: formData,
          }
        );

        if (!apiResponse.ok) {
          throw new Error("Failed to fetch results from API");
        }

        const data = await apiResponse.json();
        setResults(data);
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-4 mt-4 w-96 max-h-[600px] overflow-y-auto font-sans">
      <div className="text-center mb-4">
        <img
          src="diddo-light.png"
          alt="Diddo"
          className="mx-auto mb-2 w-[120px]"
        />
        <p className="text-sm text-gray-500">Demo</p>
      </div>
      <div className="flex justify-center mt-8">
        <button
          onClick={handleClick}
          className="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-6 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={loading}
        >
          {loading ? "Searching..." : "Find Clothes!"}
        </button>
      </div>

      {loading && (
        <div className="mt-4 text-center">
          <div className="loader mx-auto"></div>
          <p className="text-sm text-gray-500 mt-2">Finding matches...</p>
        </div>
      )}

      {error && (
        <div className="mt-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg relative" role="alert">
          <strong className="font-bold">Error:</strong>
          <span className="block sm:inline ml-2">{error}</span>
        </div>
      )}

      {results && results.items && Array.isArray(results.items) && (
        <div className="mt-6 space-y-6">
          {results.items.map((item: any) => (
            <div key={item.title}>
              <h2 className="text-lg font-semibold text-grey-100 mb-3 pb-2 border-b border-slate-200">
                {item.title}
              </h2>
              <div className="grid grid-cols-3 gap-3">
                {item.shopping_results.slice(0, 3).map((product: any) => (
                  <div
                    key={product.title}
                    className="bg-white border border-slate-200 rounded-lg p-3 shadow-sm hover:shadow-lg transition-shadow duration-200 flex flex-col justify-between"
                  >
                    <img
                      src={product.image}
                      alt={product.title}
                      className="w-full h-24 object-cover rounded-md mb-3"
                    />
                    <div className="text-center flex-grow flex flex-col">
                      <p className="font-medium text-xs text-slate-800 mb-1 flex-grow">
                        {product.title}
                      </p>
                      <p className="text-xs text-slate-500 mb-3">
                        {product.price}
                      </p>
                      <a
                        href={product.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="bg-indigo-500 hover:bg-indigo-600 text-white text-xs font-semibold px-3 py-1.5 rounded-md block w-full text-center transition-colors"
                      >
                        View
                      </a>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default App;
