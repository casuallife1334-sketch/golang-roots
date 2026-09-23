import { Component, type ReactNode } from "react";
export class ErrorBoundary extends Component<
  { children: ReactNode },
  { failed: boolean }
> {
  state = { failed: false };
  static getDerivedStateFromError() {
    return { failed: true };
  }
  render() {
    if (this.state.failed)
      return (
        <section className="center-panel" role="alert">
          <h2>Не удалось отобразить этот раздел</h2>
          <p>Данные на сервере сохранены. Попробуйте открыть раздел снова.</p>
          <button
            className="button secondary"
            onClick={() => this.setState({ failed: false })}
          >
            Повторить
          </button>
        </section>
      );
    return this.props.children;
  }
}
