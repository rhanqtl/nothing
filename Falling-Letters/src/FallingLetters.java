import java.awt.*;
import java.util.Random;

public class FallingLetters extends Component {

    private static final int DEFAULT_ROWS = 25;
    private static final int DEFAULT_COLS = 80;

    private char[][] grid;
    private int rows;
    private int cols;

    private static final char[] DEFAULT_ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ".toCharArray();

    private char[] alphabet;

    private int newLettersCountEachTime;

    private Random rand;

    public FallingLetters() {
        this(1);
    }

    public FallingLetters(final int rows, final int cols) {
        this(rows, cols, DEFAULT_ALPHABET, 1);
    }

    public FallingLetters(final char[] alphabet) {
        this(DEFAULT_ROWS, DEFAULT_COLS, alphabet, 1);
    }

    public FallingLetters(final int newLettersCountEachTime) {
        this(DEFAULT_ROWS, DEFAULT_COLS, DEFAULT_ALPHABET, newLettersCountEachTime);
    }

    public FallingLetters(final int rows, final int cols, final char[] alphabet, final int newLettersCountEachTime) {
        if (rows <= 0) {
            throw new IllegalArgumentException("Argument \"rows\" should be > 0");
        }
        if (cols <= 0) {
            throw new IllegalArgumentException("Argument \"cols\" should be > 0");
        }
        if (alphabet == null || alphabet.length == 0) {
            throw new IllegalArgumentException("Argument \"alphabet\" should not be empty");
        }
        this.rows = rows;
        this.cols = cols;
        this.grid = new char[rows][cols];
        this.alphabet = alphabet;
        this.newLettersCountEachTime = newLettersCountEachTime;
        this.rand = new Random();
        this.addKeyListener(new KeyEventHandler());
    }

    private char generateLetter() {
        return alphabet[rand.nextInt(alphabet.length)];
    }

    @Override
    public void repaint() {
        fall();
        addLetters(newLettersCountEachTime);

        for (int i = 0; i < cols * 4 + 1; i++) {
            System.out.print("=");
        }
        System.out.println();
        for (int i = 0; i < rows; i++) {
            if (i != 0) {
                for (int j = 0; j < cols * 4 + 1; j++) {
                    System.out.print("-");
                }
                System.out.println();
            }
            for (int j = 0; j < cols; j++) {
                if (j == 0) {
                    System.out.print("|");
                }
                if (isInAlphabet(grid[i][j])) {
                    System.out.print(" " + grid[i][j] + " |");
                } else {
                    System.out.print("   |");
                }
            }
            System.out.println();
        }
        for (int i = 0; i < cols * 4 + 1; i++) {
            System.out.print("=");
        }
        System.out.println();
        System.out.println();
    }

    private void addLetters(int count) {
        for (int i = 0; i < count; i++) {
            char ch = generateLetter();
            int j;
            do {
                j = rand.nextInt(cols);
            } while (isInAlphabet(grid[i][j]));
            grid[i][j] = ch;
        }
    }

    private void fall() {
        for (int i = rows - 1; i > 0; i--) {
            for (int j = 0; j < cols; j++) {
                clearCell(i, j);
                grid[i][j] = grid[i - 1][j];
                clearCell(i - 1, j);
            }
        }
    }

    private boolean isInAlphabet(char ch) {
        return ch != '\0';
    }

    private void clearCell(final int i, final int j) {
        grid[i][j] = '\0';
    }
}
