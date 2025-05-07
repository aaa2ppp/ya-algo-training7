import sys
import math
from typing import List, Tuple

class Item:
    __slots__ = ['length', 'prefix', 'middle', 'suffix']
    
    def __init__(self, length=0, prefix=0, middle=0, suffix=0):
        self.length = length
        self.prefix = prefix
        self.middle = middle
        self.suffix = suffix

def precalc(nums: List[int]) -> List[Item]:
    n = 1
    while n < len(nums):
        n *= 2
    
    tree = [Item() for _ in range(2 * n - 1)]
    
    # Заполняем листья
    for i in range(len(nums)):
        update_leaf(tree, n - 1 + i, nums[i])
    
    # Заполняем оставшиеся листья
    for i in range(n - 1 + len(nums), len(tree)):
        tree[i] = Item(length=1)
    
    # Построение дерева снизу вверх
    for i in range(n - 2, -1, -1):
        update_parent(tree, i)
    
    return tree

def update_leaf(tree: List[Item], i: int, val: int) -> None:
    if val == 0:
        tree[i] = Item(length=1, prefix=1, middle=1, suffix=1)
    else:
        tree[i] = Item(length=1)

def update_parent(tree: List[Item], i: int) -> None:
    left = 2 * i + 1
    right = 2 * i + 2
    a = tree[left]
    b = tree[right]
    
    new_item = Item(
        length=a.length + b.length,
        prefix=a.prefix,
        suffix=b.suffix,
        middle=max(a.middle, b.middle, a.suffix + b.prefix)
    )
    
    if a.prefix == a.length:
        new_item.prefix += b.prefix
    if b.suffix == b.length:
        new_item.suffix += a.suffix
    
    tree[i] = new_item

def update(tree: List[Item], i: int, val: int) -> None:
    n = (len(tree) + 1) // 2
    i += n - 1
    update_leaf(tree, i, val)
    
    while i > 0:
        i = (i - 1) // 2
        update_parent(tree, i)

def query(tree: List[Item], ql: int, qr: int) -> int:
    def dfs(i: int, l: int, r: int) -> Item:
        # Если интервал не пересекается
        if qr <= l or r <= ql:
            return Item()
        
        # Если интервал полностью покрыт
        if ql <= l and r <= qr:
            return tree[i]
        
        # Рекурсивно проверяем детей
        m = (l + r) // 2
        a = dfs(2 * i + 1, l, m)
        b = dfs(2 * i + 2, m, r)
        
        # Объединяем результаты
        res = Item(
            length=a.length + b.length,
            prefix=a.prefix,
            suffix=b.suffix,
            middle=max(a.middle, b.middle, a.suffix + b.prefix)
        )
        
        if a.prefix == a.length:
            res.prefix += b.prefix
        if b.suffix == b.length:
            res.suffix += a.suffix
        
        return res
    
    n = (len(tree) + 1) // 2
    result = dfs(0, 0, n)
    return max(result.prefix, result.middle, result.suffix)

def main():
    input = sys.stdin.read().split()
    ptr = 0
    
    n = int(input[ptr])
    ptr += 1
    
    nums = list(map(int, input[ptr:ptr + n]))
    ptr += n
    
    tree = precalc(nums)
    
    m = int(input[ptr])
    ptr += 1
    
    output = []
    for _ in range(m):
        op = input[ptr]
        ptr += 1
        
        if op == "QUERY":
            l = int(input[ptr]) - 1  # convert to 0-based
            r = int(input[ptr + 1])
            ptr += 2
            ans = query(tree, l, r)
            output.append(str(ans))
        elif op == "UPDATE":
            i = int(input[ptr]) - 1  # convert to 0-based
            v = int(input[ptr + 1])
            ptr += 2
            update(tree, i, v)
    
    print('\n'.join(output))

if __name__ == "__main__":
    # Тестовый ввод
    if len(sys.argv) > 1 and sys.argv[1] == 'test':
        test_input = """5
328 0 0 0 0
5
QUERY 1 3
UPDATE 2 832
QUERY 3 3
QUERY 2 3
UPDATE 2 0"""
        from io import StringIO
        sys.stdin = StringIO(test_input)
    
    main()