import { useCallback, useEffect, useRef } from "react";
import { useNodesInitialized, useReactFlow } from "@xyflow/react";
import { CARD_HEIGHT, CARD_WIDTH } from "./layout";

// Graph coordinates include cards that React Flow has virtualized off-screen.
export function Viewport({
  topology,
  container,
  nodes,
  focusId,
  focusRequest,
  viewportReady,
}: {
  topology: string;
  container: HTMLDivElement | null;
  nodes: { id: string; position: { x: number; y: number } }[];
  focusId?: string;
  focusRequest: number;
  viewportReady: boolean;
}) {
  const initialized = useNodesInitialized();
  const { fitView, setViewport } = useReactFlow();
  const focus = useCallback(
    (duration: number) => {
      if (!viewportReady || !container || !focusId) return;
      const node = nodes.find((item) => item.id === focusId);
      if (!node || !container.clientWidth || !container.clientHeight) return;
      const zoom = 1.25;
      void setViewport(
        {
          x: container.clientWidth / 2 - (node.position.x + CARD_WIDTH / 2) * zoom,
          y: container.clientHeight / 2 - (node.position.y + CARD_HEIGHT / 2) * zoom,
          zoom,
        },
        { duration },
      );
    },
    [container, nodes, focusId, setViewport, viewportReady],
  );
  useEffect(() => {
    focus(350);
  }, [focus, focusRequest]);
  const previous = useRef("");
  useEffect(() => {
    if (!initialized || focusId || previous.current === topology) return;
    previous.current = topology;
    void fitView({ padding: 0.2, maxZoom: 1 });
  }, [initialized, topology, fitView, focusId]);
  useEffect(() => {
    if (!container) return;
    let timer: ReturnType<typeof setTimeout>;
    let width = container.clientWidth;
    const observer = new ResizeObserver(() => {
      if (container.clientWidth === width) return;
      width = container.clientWidth;
      clearTimeout(timer);
      timer = setTimeout(() => {
        if (focusId) focus(0);
        else void fitView({ padding: 0.2, maxZoom: 1 });
      }, 120);
    });
    observer.observe(container);
    return () => {
      clearTimeout(timer);
      observer.disconnect();
    };
  }, [container, fitView, focus, focusId]);
  return null;
}
