import sys
from typing import List

def precalc(nums: List[int]) -> List[int]:
    n = 1
    while n < len(nums):
        n *= 2
    
    tree = [0] * (2 * n - 1)
    tree[n-1:n-1+len(nums)] = nums
    return tree

def update(tree: List[int], ql: int, qr: int, add: int) -> None:
    n = (len(tree) + 1) // 2
    
    def dfs(i: int, l: int, r: int) -> None:
        if qr <= l or r <= ql:
            return
        
        if ql <= l and r <= qr:
            tree[i] += add
            return
        
        m = (l + r) // 2
        dfs(2*i+1, l, m)
        dfs(2*i+2, m, r)
    
    dfs(0, 0, n)

def query(tree: List[int], idx: int) -> int:
    n = (len(tree) + 1) // 2
    target = idx + n - 1
    
    def dfs(i: int, l: int, r: int) -> int:
        if i == target:
            return tree[i]
        
        m = (l + r) // 2
        if idx < m:
            return dfs(2*i+1, l, m) + tree[i]
        else:
            return dfs(2*i+2, m, r) + tree[i]
    
    return dfs(0, 0, n)

def main():
    input = sys.stdin.read().split()
    ptr = 0
    
    n = int(input[ptr])
    ptr += 1
    
    nums = list(map(int, input[ptr:ptr+n]))
    ptr += n
    
    tree = precalc(nums)
    
    m = int(input[ptr])
    ptr += 1
    
    output = []
    for _ in range(m):
        op = input[ptr]
        ptr += 1
        
        if op == 'g':
            i = int(input[ptr]) - 1  # convert to 0-based
            ptr += 1
            ans = query(tree, i)
            output.append(str(ans))
        elif op == 'a':
            l = int(input[ptr]) - 1  # convert to 0-based
            r = int(input[ptr+1])
            v = int(input[ptr+2])
            ptr += 3
            update(tree, l, r, v)
    
    print('\n'.join(output))

if __name__ == "__main__":
    # Тестовый ввод
    if len(sys.argv) > 1 and sys.argv[1] == 'test':
        test_input = """5
2 4 3 5 2
5
g 2
g 5
a 1 3 10
g 2
g 4"""
        from io import StringIO
        sys.stdin = StringIO(test_input)
    
    main()