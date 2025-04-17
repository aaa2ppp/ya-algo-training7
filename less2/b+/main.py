def precalc(nums):
    n = len(nums)
    tab = []

    pow = [0]*(n + 1)
    j = 1

    row = [(i, nums[i]) for i in range(n)]
    tab.append(row)

    prev, k, s = row, 0, 1
    while s < len(prev):
        while j <= s*2 and j < len(pow):
            pow[j] = k
            j += 1
        row = [prev[i] if prev[i][1] > prev[i+s][1] else prev[i+s]
               for i in range(len(prev)-s)]
        tab.append(row)
        prev, k, s = row, k+1, s*2

    while j < len(pow):
        pow[j] = k
        j += 1

    return tab, pow


def query(tab, pow, l, r):
    s = r - l + 1
    k = pow[s]
    a, b = tab[k][l], tab[k][r+1-(1 << k)]
    return a[0] if a[1] > b[1] else b[0]


def main():
    n = int(input())
    nums = tuple(map(int, input().split()))
    tab, pow = precalc(nums)
    k = int(input())
    ans = []
    for _ in range(k):
        l, r = map(int, input().split())
        l, r = l-1, r-1
        idx = query(tab, pow, l, r)
        ans.append(idx+1)
    print(*ans, sep="\n")


main()
