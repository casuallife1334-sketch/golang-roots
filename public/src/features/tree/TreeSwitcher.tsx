import { Plus } from "lucide-react";
import { Select } from "../../shared/ui";
import type { Tree } from "../../types";
export function TreeSwitcher({
  trees,
  current,
  onChange,
  onCreate,
}: {
  trees: Tree[];
  current?: Tree;
  onChange: (id: string) => void;
  onCreate: () => void;
}) {
  return (
    <div className="tree-switcher">
      <Select
        aria-label="Выбрать дерево"
        value={current?.id ?? ""}
        onChange={(event) => onChange(event.target.value)}
      >
        {!current && <option value="">Выберите дерево</option>}
        {trees.map((tree) => (
          <option key={tree.id} value={tree.id}>
            {tree.name}
          </option>
        ))}
      </Select>
      <button onClick={onCreate} aria-label="Создать дерево">
        <Plus size={15} />
      </button>
    </div>
  );
}
