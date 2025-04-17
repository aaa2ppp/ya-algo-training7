def precalc(nums):
    tab = []
    row = tuple((v << 32) + i for i, v in enumerate(nums, 1))  # 1-base
    tab.append(row)
    prev, k, s = row, 0, 1
    pow = [0, 0, 0]

    while s < len(prev):
        row = tuple(prev[i] if prev[i] > prev[i+s] else prev[i+s] for i in range(len(prev)-s))
        tab.append(row)
        prev, k, s = row, k+1, s*2
        pow.extend([k]*(s*2+1-len(pow)))

    return tab, pow


def main():
    _ = int(input())
    nums = map(int, input().split())
    tab, pow = precalc(nums)
    k = int(input())
    ans = []
    for _ in range(k):
        l, r = map(int, input().split())
        k = pow[r - l + 1]
        a, b = tab[k][l-1], tab[k][r-(1 << k)]
        ans.append(a & 0xffff_ffff if a > b else b & 0xffff_ffff)
    print("\n".join(map(str, ans)))


main()
