import { useEffect, useRef } from "react";
import { useNodesInitialized, useReactFlow } from "@xyflow/react";

// Fit on topology/container changes, not on selection, photos or metadata edits.
export function Viewport({
  topology,
  container,
}: {
  topology: string;
  container: HTMLDivElement | null;
}) {
  const initialized = useNodesInitialized();
  const { fitView } = useReactFlow();
  const previous = useRef("");
  useEffect(() => {
    if (!initialized || previous.current === topology) return;
    previous.current = topology;
    void fitView({ padding: 0.2, maxZoom: 1 });
  }, [initialized, topology, fitView]);
  useEffect(() => {
    if (!container) return;
    let timer: ReturnType<typeof setTimeout>;
    let width = container.clientWidth;
    const observer = new ResizeObserver(() => {
      if (container.clientWidth === width) return;
      width = container.clientWidth;
      clearTimeout(timer);
      timer = setTimeout(() => {
        void fitView({ padding: 0.2, maxZoom: 1 });
      }, 120);
    });
    observer.observe(container);
    return () => {
      clearTimeout(timer);
      observer.disconnect();
    };
  }, [container, fitView]);
  return null;
}
