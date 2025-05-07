import sys
import os

if os.getenv("DEBUG") is not None:
    def debug(*a, **kw):
        print(*a, **kw, file=sys.stderr)
else:
    def debug(*a, **kw):
        pass


def solve(a):
    n = len(a)
    max_len = max(a).bit_length()
    bits_a = []
    bit_lens = [0]*n
    total_ones = 0
    for i in range(n):
        bits = list(map(int, str(bin(a[i]))[2:])) + \
            [0] * (max_len - a[i].bit_length())
        debug(bits, end=" -> ")
        bits.sort(reverse=True)
        bits_a.append(bits)
        ones_cnt = bits.count(1)
        total_ones += ones_cnt
        bit_lens[i] = ones_cnt
        debug(bits)

    if total_ones % 2 == 1:
        debug("oops!..")
        return None

    # Как у Гостокашина, идем от 0-го столбца, при необходимости правим столбец,
    # чтобы в нем было четное кол-во единичек. Единички будем премещать вперед,
    # но не в последнй столбец, а на *минимально* возможное растояние

    i0, j0 = 0, 0
    for i in range(max_len-1):
        debug("--------", f"i: {i}", *bits_a, sep="\n")

        bit_sum = 0
        for j in range(n):
            bit_sum += bits_a[j][i]

        if bit_sum == 0:
            # Дальше только нолики
            debug("bingo!")
            return bits_a

        if bit_sum % 2 == 0:
            debug("do nothing")
            continue

        if bit_sum > 1:
            # Перемещаем одну единичку вперед на *минимально* возможную позицию

            min_len = 100500
            min_len_j = -1
            for j in range(n):
                if bits_a[j][i] == 1 and bit_lens[j] < min_len:
                    min_len = bit_lens[j]
                    min_len_j = j

            if min_len == max_len:
                debug("oops!..")
                return None

            debug(f"swap [{min_len_j},{i}] <-> [{min_len_j},{min_len}]")
            bits_a[min_len_j][i] = 0
            bits_a[min_len_j][min_len] = 1
            bit_lens[min_len_j] += 1

        else:  # в столбце только одна единичка
            # Добавляем единичку из предыдущего столбца. Чтобы не испортить предыдущий столбец,
            # берем из него *две* единички. Вторую перемещаем вперед на *минимально*
            # возможную позицию.
            # NOTE: Мы никогда не будем здесь на первом столбце, по условию n >= 2 и
            #  a[i] >= 1

            for i0 in range(i0, i):
                j1, j2 = -1, -1
                for j in range(j0, n):
                    if bits_a[j][i] == 1:
                        continue
                    if bits_a[j][i0] == 1:
                        if j1 == -1:
                            j1 = j
                        else:
                            j2 = j
                            break

                j0 = j + 1
                if j0 >= n:
                    j0 = 0

                if j1 != -1 and j2 != -1:
                    debug(f"swap [{j1},{i0}] <-> [{j1},{i}]")
                    bits_a[j1][i0] = 0
                    bits_a[j1][i] = 1

                    debug(f"swap [{j2},{i0}] <-> [{j2},{i+1}]")
                    bits_a[j2][i0] = 0
                    bits_a[j2][i+1] = 1
                    break

            if i0 == i:
                debug("oops!..")
                return None

    debug("--------", *bits_a, sep="\n")

    bit_sum = 0
    for j in range(n):
        bit_sum ^= bits_a[j][-1]

    if bit_sum == 1:
        debug("oops!..")
        return None

    debug("bingo!")
    return bits_a


n = int(input().strip())
a = list(map(int, input().split()))

bits_a = solve(a)
if bits_a is None:
    print('impossible')
else:
    ans = []
    for i in range(n):
        # reversed здесь только для того, чтобы значения в ответе были поменьше.
        s = ''.join(map(str, reversed(bits_a[i])))
        ans.append(int(s, 2))
    print(*ans)
