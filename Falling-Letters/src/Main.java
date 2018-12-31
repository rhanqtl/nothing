public class Main {

    public static void main(String[] args) throws InterruptedException {
        FallingLetters fl = new FallingLetters(5, 5);
        for (int i = 0; i < 6; i++) {
            fl.repaint();
            Thread.sleep(1000);
        }
    }
}
