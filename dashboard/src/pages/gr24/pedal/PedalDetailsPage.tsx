import { Loader2 } from "lucide-react";
import React, { useCallback } from "react";
import { checkCredentials } from "@/lib/auth";
import { useNavigate } from "react-router-dom";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";

function PedalDetailsPage() {
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(true);

  const [socketUrl] = React.useState("ws://localhost:7001/ws/gr24/pedal");
  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

  const [jsonData, setJsonData] = React.useState([{}]);

  React.useEffect(() => {
    init();
    if (lastMessage !== null) {
      setJsonData((prev) => [
        ...prev.concat(JSON.parse(lastMessage.data)).slice(-100),
      ]);
    }
  }, [lastMessage]);

  const handleClickSendMessage = useCallback(() => sendMessage("Hello"), []);

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  const init = async () => {
    const currentRoute = window.location.pathname + window.location.search;
    const status = await checkCredentials();
    if (status != 0) {
      navigate(`/auth/register?route=${currentRoute}`);
    } else {
      setLoading(false);
    }
  };

  const LoadingComponent = () => {
    return (
      <div className="flex h-full flex-col items-center justify-center p-32">
        <Loader2 className="mr-2 h-8 w-8 animate-spin" />
      </div>
    );
  };

  return loading ? (
    <LoadingComponent />
  ) : (
    <>
      <div className="flex p-16">
        <div className="h-screen w-full text-center">
          <div>The WebSocket is currently {connectionStatus}</div>
          <Button
            className="m-2"
            onClick={handleClickSendMessage}
            disabled={readyState !== ReadyState.OPEN}
          >
            Send debug message
          </Button>
          <div className="flex w-full flex-row p-2">
            <div className="m-4 w-1/2 text-start">
              {lastMessage ? (
                <div>
                  <div>Last Message ID: {JSON.parse(lastMessage.data).id}</div>
                  <Progress
                    className="mt-4"
                    value={Math.floor(JSON.parse(lastMessage.data).apps_one)}
                  />
                  <div className="mt-2">
                    APPS One: {JSON.parse(lastMessage.data).apps_one}
                  </div>
                  <Progress
                    className="mt-4"
                    value={Math.floor(JSON.parse(lastMessage.data).apps_two)}
                  />
                  <div className="mt-2">
                    APPS Two: {JSON.parse(lastMessage.data).apps_two}
                  </div>
                  <Progress
                    className="mt-4"
                    value={Math.floor(
                      (JSON.parse(lastMessage.data).brake_pressure_front /
                        256) *
                        100,
                    )}
                  />
                  <div className="mt-2">
                    Brake Pressure Front:{" "}
                    {JSON.parse(lastMessage.data).brake_pressure_front}
                  </div>
                  <Progress
                    className="mt-4"
                    value={Math.floor(
                      (JSON.parse(lastMessage.data).brake_pressure_rear / 256) *
                        100,
                    )}
                  />
                  <div className="mt-2">
                    Brake Pressure Rear:{" "}
                    {JSON.parse(lastMessage.data).brake_pressure_rear}
                  </div>
                </div>
              ) : (
                <div>No data</div>
              )}
            </div>
            <div className="m-4 w-1/2 text-wrap text-start text-slate-500">
              <h1 className="text-xl text-white">
                Message Count: {jsonData.length}
              </h1>
              {jsonData.reverse().map((message, idx) => (
                <div key={idx}>{JSON.stringify(message)}</div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

export default PedalDetailsPage;
