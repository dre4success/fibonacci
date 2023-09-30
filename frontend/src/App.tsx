import { useEffect, useState } from "react";
import "./App.css";

const baseUrl = "http://localhost:8080/api";

function App() {
  const [current, setCurrent] = useState<number>(0);
  const [tokenHeader, setTokenHeader] = useState<string>("");
  const [error, setError] = useState<string>("");

  const fetchCurrent = async (token: string) => {
    try {
      const response = await fetch(`${baseUrl}/current`, {
        headers: {
          "X-Fib-Token": token,
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      setTokenHeader(data.token);
      setError("");
      setCurrent(data.value);
    } catch (error) {
      console.error(error);
    }
  };

  const fetchNext = async (token: string) => {
    try {
      const response = await fetch(`${baseUrl}/next`, {
        headers: { "X-Fib-Token": token, "Content-Type": "application/json" },
      });
      const data = await response.json();
      setCurrent(data.value);
    } catch (error) {
      console.error(error);
    }
  };

  const fetchPrevious = async (token: string) => {
    try {
      const response = await fetch(`${baseUrl}/previous`, {
        headers: { "X-Fib-Token": token, "Content-Type": "application/json" },
      });
      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.message);
      }
      const validData = await response.json();
      setCurrent(validData.value);
    } catch (error: any) {
      setError(error?.message);
      console.error(error);
    }
  };

  useEffect(() => {
    fetchCurrent(tokenHeader);
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Fibonacci Value</h1>
        <p>{error ? error : current}</p>
        <div>
          {error ? (
            <button onClick={() => fetchCurrent(tokenHeader)}>Current</button>
          ) : (
            <>
              <button onClick={() => fetchPrevious(tokenHeader)}>
                Previous
              </button>
              <button onClick={() => fetchNext(tokenHeader)}>Next</button>
            </>
          )}
        </div>
      </header>
    </div>
  );
}

export default App;
